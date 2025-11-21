package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		sort.Ints(b)
		g := b[0]
		for i := 1; i < n; i++ {
			g = gcd(g, b[i])
		}

		orig := make([]int, n)
		copy(orig, b)
		for i := 0; i < n; i++ {
			b[i] /= g
		}

		x := 1
		for i := 1; i < n; i++ {
			if b[i]%b[0] != 0 {
				x = gcd(x, b[i])
			}
		}
		if x == 0 {
			x = b[0]
		}
		fmt.Fprintln(out, x)
	}
}
