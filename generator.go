package sl2

import (
	"fmt"
	"hash"
	"strconv"
	"sync"

	galof "github.com/cloud9-tools/go-galoisfield"
)

const (
	DefaultGenerate256 = 3
)

type Element byte

// Generator of SL2(F_2^n) scheme
type Generator struct {
	// Generation elements
	genA ElGroup
	genB ElGroup

	// Finite field F = F_2^n
	field *galof.GF

	// inside F
	generate Element

	join, sum, multi bool

	//hash schema
	hasher hash.Hash
	lk     sync.Mutex
}

// Element of Sl2 group
type ElGroup struct {
	a, b, c, d Element
	gen        *Generator
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
	return ElGroup{a: el.a, b: el.b, c: el.c, d: el.d}
}

// Neutral element SL2 group
func Ones() ElGroup {
	return ElGroup{a: 1, b: 0, c: 0, d: 1}
}

// Generate sl2 group
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

	A := ElGroup{a: gen.generate, b: 1, c: 1, d: 0}
	B := ElGroup{a: gen.generate, b: gen.generate + 1, c: 1, d: 1}

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

func (gen *Generator) Mult(some ElGroup, other ElGroup) ElGroup {
	if !gen.Incoming(some) || !gen.Incoming(other) {
		return ElGroup{}
	}
	return ElGroup{
		a:   gen.add(gen.mul(some.a, other.a), gen.mul(some.b, other.c)),
		b:   gen.add(gen.mul(some.a, other.b), gen.mul(some.b, other.d)),
		c:   gen.add(gen.mul(some.c, other.a), gen.mul(some.d, other.c)),
		d:   gen.add(gen.mul(some.c, other.b), gen.mul(some.d, other.d)),
		gen: some.gen,
	}

}

// Build a snapshot of the data chain
func (gen *Generator) Snapshot(data ...string) ([]byte, error) {
	res := Ones()
	for _, d := range data {
		el, err := gen.Marshal(d)
		if err != nil {
			return nil, err
		}
		res = gen.Mult(res, el)
	}
	temp := res.foreach(func(in Element, temp []byte) []byte {
		return append(temp, byte(in))
	})

	return gen.hashOpts(temp, func(b []byte) ([]byte, error) {
		return temp, nil
	})
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

// Whether the element belongs to the SL2 group
func (gen *Generator) Incoming(el ElGroup) bool {
	det := gen.add(gen.mul(el.a, el.d), gen.mul(el.b, el.c))
	return det == 1
}

// String to SL2 element
func (gen *Generator) Marshal(data string) (ElGroup, error) {
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
		el = gen.Mult(el, nel)
	}
	return el, nil
}

func (gen *Generator) Snap(el ElGroup) ([]byte, error) {
	if !gen.Incoming(el) {
		return nil, fmt.Errorf("el not inside SL2 group")
	}

	temp := el.foreach(func(in Element, temp []byte) []byte {
		return append(temp, byte(in))
	})

	return gen.hashOpts(temp, func(b []byte) ([]byte, error) {
		return temp, nil
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
	gen.lk.Lock()
	defer gen.lk.Unlock()
	defer gen.hasher.Reset()

	res, err := opts(data)
	if err != nil {
		return nil, err
	}
	if _, err := gen.hasher.Write(res); err != nil {
		return nil, err
	}
	return gen.hasher.Sum(nil), nil
}
