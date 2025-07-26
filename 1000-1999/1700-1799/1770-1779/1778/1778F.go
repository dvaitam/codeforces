package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
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
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
		}
		g := a[0]
		for i := 1; i < n; i++ {
			g = gcd(g, a[i])
		}
		ans := a[0]
		if k > 0 {
			ans *= g
		}
		fmt.Fprintln(out, ans)
	}
}
