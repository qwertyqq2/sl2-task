package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/qwertyqq2/sl2"
)

type Stats struct {
	oper int
	time time.Duration
}

var (
	filename = "testing64big.txt"
	mu       sync.Mutex
	baseSize = 30
)

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

func generatorFactory() *sl2.Generator {
	return sl2.Generate(
		sl2.SetDefaultElement(),
		sl2.SetSha256(),
		sl2.SetOrderField64(),
	)
}

type finder struct {
	ctx    context.Context
	cancel context.CancelFunc

	base []byte

	gen  *sl2.Generator
	oper int
	file *os.File

	ch chan *Stats

	lastUpdate time.Time
}

func newFinder(ctx context.Context, file *os.File, base []byte) *finder {
	gen := generatorFactory()
	ctx, cancel := context.WithCancel(ctx)
	return &finder{
		ctx:    ctx,
		cancel: cancel,
		gen:    gen,
		ch:     make(chan *Stats, 50),
		file:   file,
		base:   base,
	}
}

func (finder *finder) find(countFinder int) {
	if countFinder <= 0 {
		return
	}
	for i := 0; i < countFinder; i++ {
		go find(finder.ctx, finder.cancel, finder.base, func() ([]byte, error) {
			str := randStringRunes(baseSize)
			return finder.gen.Snapshot(str)
		}, func(stats *Stats) {
			select {
			case <-finder.ctx.Done():
				return
			case finder.ch <- stats:
			}
		})
	}
}

func (finder *finder) run(countFinder int) {
	go finder.find(countFinder)
	for {
		select {
		case stat := <-finder.ch:
			mu.Lock()
			fmt.Println(stat.str())
			if _, err := finder.file.Write([]byte(stat.str())); err != nil {
				log.Fatal(err)
			}
			mu.Unlock()

		case <-finder.ctx.Done():
			return
		}
	}
}

func main() {
	gen := generatorFactory()

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	base, err := gen.Snapshot(randStringRunes(baseSize))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)

	var countFinder = 3

	defer cancel()

	for i := 0; i < 3; i++ {
		go newFinder(ctx, f, base).run(countFinder)
	}

	select {
	case <-ctx.Done():
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
				hand(&Stats{oper: counter, time: time.Duration(t.Seconds())})
				lastUpdate = time.Now()
				counter = 0
			}
			counter++
		}
	}
}
