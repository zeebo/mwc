package mwc

import "testing"

func run[V any](fn func() V) func(b *testing.B) {
	return func(b *testing.B) {
		b.Run("Sequential", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = fn()
			}
		})
		b.Run("Parallel", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_ = fn()
				}
			})
		})
	}
}

func runN[V any](fn func(V) V) func(b *testing.B) {
	return run(func() V { return fn(*new(V)) })
}

func BenchmarkGlobal(b *testing.B) {
	b.Run("Uint64", run(Uint64))
	b.Run("Uint64n", runN(Uint64n))
	b.Run("Uint32", run(Uint32))
	b.Run("Uint32n", runN(Uint32n))
	b.Run("Intn", runN(Intn))
	b.Run("Float64", run(Float64))
	b.Run("Float32", run(Float32))
}
