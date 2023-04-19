package sl2

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

const (
	DefualtGenEl = 4
)

func TestAssoc(t *testing.T) {
	gen := Generate(
		SetDefaultElement(),
		SetSha256(),
	)

	h1, err := gen.Marshal("asfasf")
	if err != nil {
		t.Fatal(err)
	}
	h2, err := gen.Marshal("qwe")
	if err != nil {
		t.Fatal(err)
	}

	h3, err := gen.Marshal("onskdmpldmpfe")
	if err != nil {
		t.Fatal(err)
	}

	if !gen.Incoming(h1) || !gen.Incoming(h2) || !gen.Incoming(h3) {
		t.Fatal()
	}

	chain1 := []string{"block1, block2, block3, block4, block5"}

	snap1, err := gen.Snapshot(chain1...)
	if err != nil {
		t.Fatal(err)
	}

	chain2 := []string{"block1, block3, block2, block4, block5"}

	snap2, err := gen.Snapshot(chain2...)
	if err != nil {
		t.Fatal(err)
	}

	chain3 := []string{"block5, block2, block3, block4, block1"}

	snap3, err := gen.Snapshot(chain3...)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(snap1, snap2) || bytes.Equal(snap2, snap3) {
		t.Fatal()
	}

	chainContintie := []string{"block6, block7, block8"}
	chain := append(chain1, chainContintie...)

	snapChain, err := gen.Snapshot(chain...)
	if err != nil {
		t.Fatal(err)
	}

	idx := int(len(chain) / 2)
	base := chain[idx]

	snap, err := bildChainSnapshot(t, idx, base, chain, gen)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(snap, snapChain) {
		t.Fatal()
	}

	fmt.Println("the associativity property is satisfied")
}

func bildChainSnapshot(t *testing.T, idx int, base string, chain []string, gen *Generator) ([]byte, error) {
	chainBefore := chain[:idx]
	chainAfter := chain[idx+1:]

	snapBefore, err := gen.SnapshotElement(chainBefore...)
	if err != nil {
		return nil, err
	}
	snapAfter, err := gen.SnapshotElement(chainAfter...)
	if err != nil {
		return nil, err
	}

	elBase, err := gen.Marshal(base)
	if err != nil {
		return nil, err
	}

	m, err := gen.Mult(snapBefore, elBase)
	if err != nil {
		return nil, err
	}

	res, err := gen.Mult(m, snapAfter)
	if err != nil {
		return nil, err
	}

	return res.Snap()

}

func TestMarshal(t *testing.T) {
	data := "qwedsgreregs"
	bytes, err := marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(bytes)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestEqualInput(t *testing.T) {
	gen := Generate(
		SetDefaultElement(),
		SetSha256(),
	)

	h1, err := gen.Snapshot("qweqwe")
	if err != nil {
		t.Fatal(err)
	}
	h2, err := gen.Snapshot("qweqwe")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(h1, h2) {
		t.Fatal()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var (
		equals, nequals = 0, 0
	)

	base := randStringRunes(10)
	baseHash, err := gen.Snapshot(base)
	if err != nil {
		t.Fatal(err)
	}

	for {
		select {
		case <-ctx.Done():
			goto Finish
		default:
		}
		str := randStringRunes(10)
		hash, err := gen.Snapshot(str)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Equal(baseHash, hash) {
			equals++
		} else {
			nequals++
		}
	}

Finish:

	fmt.Println("equals", equals)
	fmt.Println("nequals", nequals)
}

func TestAsHash(t *testing.T) {
	gen := Generate(
		SetOrderField128(),
		SetDefaultElement(),
		SetSha256(),
	)

	data := "my name is shao khan"
	strs := strings.Split(data, " ")

	hash, err := gen.Snapshot(strs...)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(hash))
}
