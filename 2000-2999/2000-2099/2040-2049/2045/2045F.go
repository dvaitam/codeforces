package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

// For a cell in row r (1-indexed from top), define e = (N - r) % (K+1).
// A single stone at that row has Grundy number g = 1 << e.
// If there are at least K rows below (r + K <= N), the K descendants cover
// all bits except e, so the span of reachable nimbers misses only bit e.
// In that case, for a pile with t = stones % (K+1):
//
//	SG = 0 if t == 0
//	SG = (t/2) * 2^(K+1) ^ ((t&1) * g) otherwise (bits are disjoint).
//
// If fewer than K rows are below, the reachable span only uses lower bits
// and the pile behaves like a normal subtraction game scaled by g, i.e.
//
//	SG = (t * g) with period K+1.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var N, M int64
		var K int64
		fmt.Fscan(in, &N, &M, &K)
		mod := K + 1
		bigShift := uint(K + 1)

		var nim big.Int
		var tmp big.Int

		for i := int64(0); i < M; i++ {
			var r, c, a int64
			fmt.Fscan(in, &r, &c, &a)
			t := a % mod
			if t == 0 {
				continue
			}
			e := uint((N - r) % mod)
			if N-r < K { // not enough rows to use full window
				tmp.SetInt64(t)
				tmp.Lsh(&tmp, e)
				nim.Xor(&nim, &tmp)
			} else { // full K rows available
				if t&1 == 1 {
					tmp.SetInt64(1)
					tmp.Lsh(&tmp, e)
					nim.Xor(&nim, &tmp)
				}
				if high := t / 2; high > 0 {
					tmp.SetInt64(high)
					tmp.Lsh(&tmp, bigShift)
					nim.Xor(&nim, &tmp)
				}
			}
		}

		if nim.Sign() != 0 {
			fmt.Fprintln(out, "Anda")
		} else {
			fmt.Fprintln(out, "Kamu")
		}
	}
}
