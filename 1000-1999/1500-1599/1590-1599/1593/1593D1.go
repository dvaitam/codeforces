package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		allEqual := true
		for i := 1; i < n; i++ {
			if a[i] != a[0] {
				allEqual = false
				break
			}
		}
		if allEqual {
			fmt.Fprintln(out, -1)
			continue
		}
		g := 0
		for i := 1; i < n; i++ {
			d := abs(a[i] - a[0])
			g = gcd(g, d)
		}
		fmt.Fprintln(out, g)
	}
}
