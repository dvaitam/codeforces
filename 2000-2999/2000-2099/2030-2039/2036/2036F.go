package main

import (
	"bufio"
	"fmt"
	"os"
)

func xorPrefix(n int64) int64 {
	switch n & 3 {
	case 0:
		return n
	case 1:
		return 1
	case 2:
		return n + 1
	default:
		return 0
	}
}

func xorRange(a, b int64) int64 {
	if a > b {
		return 0
	}
	return xorPrefix(b) ^ xorPrefix(a-1)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var l, r int64
		var i int
		var k int64
		fmt.Fscan(reader, &l, &r, &i, &k)

		total := xorRange(l, r)

		if i == 0 {
			// modulus 1, all numbers excluded if k=0
			fmt.Fprintln(writer, 0)
			continue
		}

		m := int64(1) << uint(i)

		k %= m
		delta := (k - l) % m
		if delta < 0 {
			delta += m
		}
		first := l + delta
		if first > r {
			fmt.Fprintln(writer, total)
			continue
		}
		cnt := (r-first)/m + 1
		start := (first - k) / m
		xorT := xorRange(start, start+cnt-1)
		var subset int64
		if cnt&1 == 1 {
			subset = k
		}
		subset ^= xorT << uint(i)
		fmt.Fprintln(writer, total^subset)
	}
}
