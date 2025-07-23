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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	ones := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		if a[i] == 1 {
			ones++
		}
	}

	if ones > 0 {
		fmt.Fprintln(out, n-ones)
		return
	}

	g := a[0]
	for i := 1; i < n; i++ {
		g = gcd(g, a[i])
	}
	if g != 1 {
		fmt.Fprintln(out, -1)
		return
	}

	best := n + 1
	for i := 0; i < n; i++ {
		cur := 0
		for j := i; j < n; j++ {
			if cur == 1 {
				break
			}
			cur = gcd(cur, a[j])
			if cur == 1 {
				if j-i+1 < best {
					best = j - i + 1
				}
				break
			}
		}
	}

	ans := best - 1 + n - 1
	fmt.Fprintln(out, ans)
}
