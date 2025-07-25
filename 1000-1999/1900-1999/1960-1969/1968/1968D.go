package main

import (
	"bufio"
	"fmt"
	"os"
)

func bestScore(n int, k int64, start int, p []int, a []int64) int64 {
	pos := start - 1
	best := int64(k) * a[pos]
	prefix := int64(0)
	limit := int64(n)
	if k < limit {
		limit = k
	}
	for i := int64(1); i < limit; i++ {
		prefix += a[pos]
		pos = p[pos] - 1
		cand := prefix + (int64(k)-i)*a[pos]
		if cand > best {
			best = cand
		}
	}
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var k int64
		var PB, PS int
		fmt.Fscan(reader, &n, &k, &PB, &PS)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &p[i])
		}
		aVals := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &aVals[i])
		}
		sb := bestScore(n, k, PB, p, aVals)
		ss := bestScore(n, k, PS, p, aVals)
		if sb > ss {
			fmt.Fprintln(writer, "Bodya")
		} else if ss > sb {
			fmt.Fprintln(writer, "Sasha")
		} else {
			fmt.Fprintln(writer, "Draw")
		}
	}
}
