package atomicutil

import (
	"math"
	"sync/atomic"
)

// AppendMean atomically appends a value to a float32 mean, float32 value
// and uint32 counter to calculate a running mean, and return the current
// running mean
func AppendMean(ptr *uint64, n float32) float32 {
	for {
		old := atomic.LoadUint64(ptr)
		oVBits := old >> 32
		oN := uint32(old)
		oV := math.Float32frombits(uint32(oVBits))
		oV += n
		nVBits := math.Float32bits(oV)
		oN++
		new := uint64(oN) | uint64(nVBits)<<32
		if atomic.CompareAndSwapUint64(ptr, old, new) {
			return oV
		}
	}
}

// Tx performs an atomic transaction on a uint64, using a compare and swap
// operation. the transaction should be as light as possible, as when the
// CAS operation fails we burn CPU until we can succeed. ideally this would
// only be used to implement your own bit packing with a uint64, rather than
// treating this like holding a mutex on the value
func Tx(ptr *uint64, fn func(uint64) (uint64, error)) (uint64, error) {
	for {
		old := atomic.LoadUint64(ptr)
		new, err := fn(old)
		if err != nil {
			return new, err
		}
		if new == old {
			return new, nil
		}
		if atomic.CompareAndSwapUint64(ptr, old, new) {
			return new, nil
		}
	}
}
