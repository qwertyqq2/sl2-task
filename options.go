package sl2

import (
	"crypto/md5"
	"crypto/sha256"

	galof "github.com/cloud9-tools/go-galoisfield"
)

type Option func(gen *Generator)

// Default element of the final field generating SL2 group
func SetDefaultElement() Option {
	return func(gen *Generator) {
		gen.generate = Element(2)
	}
}

func SetOrderField256() Option {
	return func(gen *Generator) {
		gen.field = galof.DefaultGF256
	}
}

func SetOrderField128() Option {
	return func(gen *Generator) {
		gen.field = galof.DefaultGF128
	}
}

func SetOrderField64() Option {
	return func(gen *Generator) {
		gen.field = galof.DefaultGF64
	}
}

func SetOrderField32() Option {
	return func(gen *Generator) {
		gen.field = galof.DefaultGF32
	}
}

func SetJoinBytesForHash(join bool) Option {
	return func(gen *Generator) {
		gen.join = join
	}
}

func SetSumBytesForHash(sum bool) Option {
	return func(gen *Generator) {
		gen.sum = sum
	}
}

func SetMulBytesForHash(multi bool) Option {
	return func(gen *Generator) {
		gen.multi = multi
	}
}

func SetSha256() Option {
	return func(gen *Generator) {
		gen.hasher = sha256.New()
	}
}

func SetMd5() Option {
	return func(gen *Generator) {
		gen.hasher = md5.New()
	}
}

func defaultOpts() Option {
	return func(gen *Generator) {
		gen.join = false
		gen.sum = false
		gen.multi = false
		gen.hasher = sha256.New()
		gen.field = galof.DefaultGF128
	}
}
