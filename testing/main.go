package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/qwertyqq2/sl2"
)

type Stats struct {
	oper int
	time time.Duration
}

func (s *Stats) str() string {
	return fmt.Sprintf("oper: %d, timedur: %d\n", s.oper, s.time)
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

func main() {
	gen := sl2.Generate256Defualt(
		sl2.SetDefaultElement(),
		sl2.SetSha256(),
	)

	base, err := gen.Snap(randStringRunes(10))
	if err != nil {
		log.Fatal(err)
	}

	statsCh := make(chan *Stats, 5)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)

	var countFinder = 10

	defer cancel()

	for i := 0; i < countFinder; i++ {
		go find(
			ctx, cancel, base,
			func() ([]byte, error) {
				str := randStringRunes(10)
				return gen.Snap(str)
			}, func(stats *Stats) {
				select {
				case <-ctx.Done():
					return
				case statsCh <- stats:

				}
			})
	}

	for {
		select {
		case stat := <-statsCh:
			fmt.Println(stat.str())
		case <-ctx.Done():
			return
		}
	}
}

func find(
	ctx context.Context,
	stop context.CancelFunc,
	base []byte,
	bild func() ([]byte, error),
	hand func(stats *Stats)) {
	defer stop()

	var (
		counter    int
		lastUpdate time.Time
	)

	lastUpdate = time.Now()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			snap, err := bild()
			if err != nil {
				return
			}
			if bytes.Equal(snap, base) {
				t := time.Now().Sub(lastUpdate)
				hand(&Stats{oper: counter, time: t})
				lastUpdate = time.Now()
			}
			counter++
		}
	}
}
