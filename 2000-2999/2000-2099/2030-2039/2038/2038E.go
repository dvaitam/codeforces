package main

import (
	"bufio"
	"fmt"
	"os"
)

func nextInt(r *bufio.Reader) int64 {
	var x int64
	sign := int64(1)
	c, _ := r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		x = x*10 + int64(c-'0')
		c, _ = r.ReadByte()
	}
	return x * sign
}

func main() {
	in := bufio.NewReader(os.Stdin)

	n := int(nextInt(in))
	v := make([]int64, n)
	for i := 0; i < n; i++ {
		v[i] = nextInt(in)
	}
	h := make([]int64, n-1)
	for i := 0; i+1 < n; i++ {
		h[i] = nextInt(in)
	}

	// cur holds the maximum water that can end up in barrels [2..k] after
	// processing pipes to the right of position k. Start with all water except
	// the first barrel.
	cur := float64(0)
	for i := 1; i < n; i++ {
		cur += float64(v[i])
	}

	// Process pipes from right to left, excluding the first pipe (between barrels 1 and 2).
	for i := n - 2; i >= 1; i-- {
		s := float64(i) // number of barrels on the left side within the subchain [2..i+1]
		hi := float64(h[i])
		if cur >= s*hi {
			cur = s * (cur + hi) / (s + 1.0)
		}
		// Otherwise, the rightmost barrel can be emptied completely, so cur stays the same.
	}

	// Finally, connect the first barrel via the first pipe.
	total := float64(v[0]) + cur
	ans := total
	if total >= float64(h[0]) {
		ans = (total + float64(h[0])) / 2.0
	}

	fmt.Printf("%.15f\n", ans)
}
