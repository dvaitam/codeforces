package main

import (
	"bufio"
	"fmt"
	"os"
)

func charID(c byte) int {
	switch c {
	case 'R':
		return 0
	case 'P':
		return 1
	default:
		return 2
	}
}

func delta(a, b int) int {
	if a == b {
		return -1
	}
	if b == (a+1)%3 {
		return 0
	}
	return 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		n := len(s)
		vals := make([]int, n)
		for i := 0; i < n; i++ {
			vals[i] = charID(s[i])
		}

		diffS := 0
		for i := 0; i < n-1; i++ {
			diffS += delta(vals[i], vals[i+1])
		}

		best := int(1e9)

		try := func(pLen int, lastChar int, diffPref int) {
			if pLen == 0 {
				if vals[0] != 0 {
					return
				}
			}
			diff := diffPref + diffS
			if pLen > 0 {
				diff += delta(lastChar, vals[0])
			}
			extra := 0
			if diff < 0 {
				extra = -diff
			}
			cand := n + pLen + extra
			if cand < best {
				best = cand
			}
		}

		try(0, 0, 0) // substring starts immediately (needs s[0] == 'R')
		try(1, 0, 0) // prefix "R"
		try(2, 2, 1) // prefix "RS"
		try(3, 1, 2) // prefix "RSP"

		fmt.Fprintln(out, best)
	}
}
