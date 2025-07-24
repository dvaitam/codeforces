package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	s := make([]int, n)
	g := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &s[i], &g[i])
	}

	lo := make([]int, n)
	hi := make([]int, n)
	for i := 0; i < n; i++ {
		lo[i] = s[i]
		hi[i] = s[i] + g[i]
	}

	for i := 1; i < n; i++ {
		if lo[i] < lo[i-1]-1 {
			lo[i] = lo[i-1] - 1
		}
		if hi[i] > hi[i-1]+1 {
			hi[i] = hi[i-1] + 1
		}
	}
	for i := n - 2; i >= 0; i-- {
		if lo[i] < lo[i+1]-1 {
			lo[i] = lo[i+1] - 1
		}
		if hi[i] > hi[i+1]+1 {
			hi[i] = hi[i+1] + 1
		}
	}
	for i := 0; i < n; i++ {
		if lo[i] > hi[i] || hi[i] < s[i] {
			fmt.Fprintln(out, -1)
			return
		}
	}

	total := int64(0)
	for i := 0; i < n; i++ {
		total += int64(hi[i] - s[i])
	}
	fmt.Fprintln(out, total)
	for i := 0; i < n; i++ {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, hi[i])
	}
	out.WriteByte('\n')
}
