package mwc

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	rngState uint64
	rngInc   uint64 = uint64(time.Now().UnixNano()) | 1
)

func Rand() *T { return New(atomic.AddUint64(&rngState, rngInc), rngInc) }

var mwcPool = sync.Pool{New: func() interface{} { return Rand() }}

func with[V any](f func(*T) V) V {
	p, _ := mwcPool.Get().(*T)
	v := f(p)
	mwcPool.Put(p)
	return v
}

func withN[V any](n V, f func(*T, V) V) V {
	p, _ := mwcPool.Get().(*T)
	v := f(p, n)
	mwcPool.Put(p)
	return v
}

func Uint64() uint64          { return with((*T).Uint64) }
func Uint64n(n uint64) uint64 { return withN(n, (*T).Uint64n) }
func Uint32() uint32          { return with((*T).Uint32) }
func Uint32n(n uint32) uint32 { return withN(n, (*T).Uint32n) }
func Intn(n int) int          { return withN(n, (*T).Intn) }
func Float64() float64        { return with((*T).Float64) }
func Float32() float32        { return with((*T).Float32) }
