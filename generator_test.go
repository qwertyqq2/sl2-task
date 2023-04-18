package sl2

import (
	"fmt"
	"testing"
)

const (
	DefualtGenEl = 4
)

func TestMarshal(t *testing.T) {
	data := "qwedsgreregs"
	bytes, err := marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(bytes)
}

func TestHash(t *testing.T) {
	gen := Generate256Defualt(
		SetDefaultElement(5),
		SetSumBytesForHash(false),
		SetSha256(),
	)

	name := "qweqweqw"
	res, err := gen.Snap(name)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}
