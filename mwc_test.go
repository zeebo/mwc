package mwc

import (
	"runtime"
	"testing"

	"github.com/zeebo/assert"
)

func BenchmarkRNG_Next(b *testing.B) {
	var hole uint64
	r := New(0, 0)
	b.SetBytes(8)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hole += r.Uint64()
	}

	runtime.KeepAlive(hole)
}

func BenchmarkNew(b *testing.B) {
	var r T

	for i := 0; i < b.N; i++ {
		r = New(0, 0)
	}

	runtime.KeepAlive(r)
}

func BenchmarkRNG_Uint64n(b *testing.B) {
	var hole uint64
	r := New(0, 0)
	b.SetBytes(8)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hole += r.Uint64n(1000)
	}

	runtime.KeepAlive(hole)
}

func BenchmarkRNG_Uint32n(b *testing.B) {
	var hole uint32
	r := New(0, 0)
	b.SetBytes(4)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hole += r.Uint32n(1000)
	}

	runtime.KeepAlive(hole)
}

func BenchmarkRNG_Float64(b *testing.B) {
	var hole float64
	r := New(0, 0)
	b.SetBytes(8)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hole += r.Float64()
	}

	runtime.KeepAlive(hole)
}

func TestRNG_Known(t *testing.T) {
	r := New(0xb01df00ddeadbeef, 0xcafefade1337d00d)
	assert.Equal(t, r.Uint64(), uint64(0xdfaa30de6f67341e))
	assert.Equal(t, r.Uint64(), uint64(0x7d6ca66d7da03a73))
	assert.Equal(t, r.Uint64(), uint64(0xc84f6fea88c4ae40))
	assert.Equal(t, r.Uint64(), uint64(0x6c9edfd9edfd2f09))
	assert.Equal(t, r.Uint64(), uint64(0x66497b2025e54253))
	assert.Equal(t, r.Uint64(), uint64(0x032137037bdcbf89))
}

func TestRNG_Uint64n(t *testing.T) {
	r := New(1, 2)
	for i := 0; i < 1000000; i++ {
		m := r.Uint64()
		assert.That(t, r.Uint64n(m) < m)
		assert.That(t, r.Uint64n(10) < 10)
	}
}

func TestRNG_Uint32n(t *testing.T) {
	r := New(3, 4)
	for i := 0; i < 1000000; i++ {
		m := r.Uint32()
		assert.That(t, r.Uint32n(m) < m)
		assert.That(t, r.Uint32n(10) < 10)
	}
}
