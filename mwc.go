package mwc

import (
	"encoding/binary"
	"math/bits"
	"sync/atomic"
	"time"
)

type T struct {
	x1 uint64
	x2 uint64
	x3 uint64
	c  uint64
}

const m64 = 0xfeb344657c0af413

var (
	rngState uint64
	rngInc   uint64 = uint64(time.Now().UnixNano())
)

func Rand() T { return New(atomic.AddUint64(&rngState, rngInc), rngInc) }

func New(k1, k2 uint64) (r T) {
	const (
		k = 0xcafef00dd15ea5e5
		c = 0x14057b7ef767814f

		rs = (((k * m64) & (1<<64 - 1)) + c) & (1<<64 - 1)
		rc = 0

		m1l uint64 = (k * m64) & (1<<64 - 1)
		m1h uint64 = (k * m64) >> 64

		m4l uint64 = (rs * m64) & (1<<64 - 1)
		m4h uint64 = (rs * m64) >> 64
	)

	var b uint64

	m2h, m2l := bits.Mul64(k2, m64)
	m3h, m3l := bits.Mul64(k1, m64)

	r.x2, b = bits.Add64(m2l, m1h, rc)
	r.x1, b = bits.Add64(m3l, m2h, b)

	m5h, m5l := bits.Mul64(r.x2, m64)
	m6h, m6l := bits.Mul64(r.x1, m64)

	r.x3, b = bits.Add64(m4l, m3h, b)
	r.x2, b = bits.Add64(m5l, m4h, b)
	r.x1, b = bits.Add64(m6l, m5h, b)
	r.c = m6h + b

	return
}

func (r *T) Int63() int64    { return int64(r.Uint64() & (1<<63 - 1)) }
func (r *T) Seed(seed int64) { *r = New(uint64(seed), uint64(seed)) }

func (r *T) Uint64() uint64 {
	h, l := bits.Mul64(r.x3, m64)
	o := (r.x3 ^ r.x2) + (r.x1 ^ h)
	x1, b := bits.Add64(l, r.c, 0)
	*r = T{x1, r.x1, r.x2, h + b}
	return o
}

func (r *T) Uint64n(n uint64) uint64 {
	if n == 0 {
		return 0
	}

	x := r.Uint64()
	h, l := bits.Mul64(x, n)

	if l < n {
		t := -n
		if t >= n {
			t -= n
			if t >= n {
				t = t % n
			}
		}

		for l < t {
			x = r.Uint64()
			h, l = bits.Mul64(x, n)
		}
	}

	return h
}

func (r *T) Uint32() uint32          { return uint32(r.Uint64()) }
func (r *T) Uint32n(n uint32) uint32 { return uint32(r.Uint64n(uint64(n))) }

func (r *T) Intn(n int) int {
	if n < 0 {
		return 0
	}
	return int(r.Uint64n(uint64(n)))
}

func (r *T) Float64() (v float64) {
	for {
		v = float64(r.Uint64()>>11) / (1 << 53)
		if v != 1 {
			return
		}
	}
}

func (r *T) Float32() (v float32) {
	for {
		v = float32(r.Uint32()>>8) / (1 << 24)
		if v != 1 {
			return
		}
	}
}

func (r *T) Read(p []byte) (n int, err error) {
	n = len(p)
	for len(p) > 8 {
		binary.LittleEndian.PutUint64(p[:8], r.Uint64())
		p = p[8:]
	}
	x := r.Uint64()
	for len(p) > 0 {
		p[0] = byte(x)
		x >>= 8
		p = p[1:]
	}
	return
}
