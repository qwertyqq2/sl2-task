package sl2

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const (
	DefualtGenEl = 4
)

func TestAssoc(t *testing.T) {
	gen := Generate256Defualt(
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

func TestHash(t *testing.T) {
	gen := Generate256Defualt(
		SetDefaultElement(),
		SetSha256(),
	)

	h1, err := gen.Snapshot("qweqwe")
	if err != nil {
		t.Fatal(err)
	}
	h2, err := gen.Snapshot("qwe")
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
