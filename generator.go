package sl2

import (
	"fmt"
	"hash"
	"strconv"

	galof "github.com/cloud9-tools/go-galoisfield"
)

const (
	DefaultGenerate256 = 5
)

type Element byte

type Generator struct {
	genA  ElGroup
	genB  ElGroup
	field *galof.GF

	generate Element

	join, sum, multi bool
	hasher           hash.Hash
}

type ElGroup struct {
	a, b, c, d Element
}

func (el ElGroup) foreach(fn func(in Element, temp []byte) []byte) []byte {
	var res []byte
	res = fn(el.a, res)
	res = fn(el.b, res)
	res = fn(el.c, res)
	res = fn(el.d, res)

	return res
}

func copy(el ElGroup) ElGroup {
	return ElGroup{
		a: el.a,
		b: el.b,
		c: el.c,
		d: el.d,
	}
}

func Generate256Defualt(opts ...Option) *Generator {
	field := galof.DefaultGF256
	gen := &Generator{
		field: field,
	}

	defaultOpts()(gen)
	for _, o := range opts {
		o(gen)
	}

	if gen.generate == 0 {
		gen.generate = DefaultGenerate256
	}

	A := ElGroup{
		a: gen.generate,
		b: 1,
		c: 1,
		d: 0,
	}

	B := ElGroup{
		a: gen.generate,
		b: gen.generate + 1,
		c: 1,
		d: 1,
	}
	gen.genA = A
	gen.genB = B

	return gen
}

func (gen *Generator) mul(a, b Element) byte {
	return gen.field.Mul(byte(a), byte(b))
}

func (gen *Generator) add(a, b byte) Element {
	return Element(gen.field.Add(byte(a), byte(b)))
}

func (gen *Generator) mult(some ElGroup, other ElGroup) ElGroup {
	return ElGroup{
		a: gen.add(gen.mul(some.a, other.a), gen.mul(some.b, other.c)),
		b: gen.add(gen.mul(some.a, other.b), gen.mul(some.b, other.d)),
		c: gen.add(gen.mul(some.c, other.a), gen.mul(some.d, other.c)),
		d: gen.add(gen.mul(some.c, other.b), gen.mul(some.d, other.d)),
	}

}

func marshal(data string) ([]int, error) {
	var b []byte
	var res []int
	for _, c := range data {
		b = strconv.AppendInt(b, int64(c), 2)
	}
	for _, bb := range b {
		switch bb {
		case 49:
			res = append(res, 1)
		case 48:
			res = append(res, 0)
		default:
			return nil, fmt.Errorf("undifined byte")
		}
	}
	return res, nil
}

func (gen *Generator) snap(data string) (ElGroup, error) {
	bin, err := marshal(data)
	if err != nil {
		return ElGroup{}, err
	}

	b := bin[0]
	el, err := gen.pi(b)
	if err != nil {
		return ElGroup{}, err
	}
	for i := 1; i < len(bin); i++ {
		b := bin[i]
		nel, err := gen.pi(b)
		if err != nil {
			return ElGroup{}, err
		}
		el = gen.mult(el, nel)
	}
	return el, nil
}

func (gen *Generator) Snap(data string) ([]byte, error) {
	el, err := gen.snap(data)
	if err != nil {
		return nil, err
	}

	temp := el.foreach(func(in Element, temp []byte) []byte {
		return append(temp, byte(in))
	})

	if gen.join {

	}

	if gen.sum {

	}

	if gen.multi {

	}

	return gen.hashOpts(temp, func(b []byte) ([]byte, error) {
		return b, nil
	})
}

func (gen *Generator) pi(v int) (ElGroup, error) {
	switch v {
	case 1:
		return gen.genA, nil
	case 0:
		return gen.genB, nil
	default:
		return ElGroup{}, fmt.Errorf("undefined byte")
	}
}

func (gen *Generator) hashOpts(data []byte, opts func([]byte) ([]byte, error)) ([]byte, error) {
	res, err := opts(data)
	if err != nil {
		return nil, err
	}
	if _, err := gen.hasher.Write(res); err != nil {
		return nil, err
	}
	return gen.hasher.Sum(nil), nil
}
