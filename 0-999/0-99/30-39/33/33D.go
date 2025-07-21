package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	xs := make([]int64, n)
	ys := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &xs[i], &ys[i])
	}
	rs := make([]int64, m)
	cx := make([]int64, m)
	cy := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &rs[i], &cx[i], &cy[i])
	}
	// bitsets of circles containing each point
	wlen := (m + 63) / 64
	bs := make([][]uint64, n)
	for i := 0; i < n; i++ {
		b := make([]uint64, wlen)
		xi, yi := xs[i], ys[i]
		for j := 0; j < m; j++ {
			dx := xi - cx[j]
			dy := yi - cy[j]
			if dx*dx+dy*dy < rs[j]*rs[j] {
				idx := j >> 6
				off := uint(j & 63)
				b[idx] |= 1 << off
			}
		}
		bs[i] = b
	}
	// answer queries
	for qi := 0; qi < k; qi++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		a--
		b--
		cnt := 0
		ba, bb := bs[a], bs[b]
		for i := 0; i < wlen; i++ {
			cnt += bits.OnesCount64(ba[i] ^ bb[i])
		}
		fmt.Fprintln(writer, cnt)
	}
}
