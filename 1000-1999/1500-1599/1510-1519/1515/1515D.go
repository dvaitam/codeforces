package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve(reader, writer)
	}
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n, l, r int
	fmt.Fscan(reader, &n, &l, &r)
	colors := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &colors[i])
	}

	left := make([]int, n+1)
	right := make([]int, n+1)
	for i := 0; i < l; i++ {
		left[colors[i]]++
	}
	for i := l; i < n; i++ {
		right[colors[i]]++
	}

	// pair socks of the same color
	for c := 1; c <= n; c++ {
		m := left[c]
		if right[c] < m {
			m = right[c]
		}
		left[c] -= m
		right[c] -= m
		l -= m
		r -= m
	}

	// ensure left side has at least as many socks as right
	if l < r {
		left, right = right, left
		l, r = r, l
	}
	diff := l - r

	// count how many pairs can be formed inside the larger side
	pairs := 0
	for c := 1; c <= n; c++ {
		pairs += left[c] / 2
	}

	use := diff / 2
	if pairs < use {
		use = pairs
	}
	res := use
	diff -= use * 2
	l -= use * 2

	// move remaining socks to balance sides
	res += diff / 2
	l -= diff / 2
	r += diff / 2

	// remaining unmatched socks need recoloring
	res += l

	fmt.Fprintln(writer, res)
}
