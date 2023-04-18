package sl2

import (
	"crypto/md5"
	"crypto/sha256"
)

type Option func(gen *Generator)

func SetDefaultElement(el byte) Option {
	return func(gen *Generator) {
		gen.generate = Element(el)
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
	}
}
