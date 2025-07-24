package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	seen := make([]bool, 1001)
	less := make([]int, 0)
	greater := make([]int, 0)
	minA, maxA := 1001, -1

	for i := 0; i < k; i++ {
		var a int
		fmt.Fscan(in, &a)
		if seen[a] {
			continue
		}
		seen[a] = true
		if a == n {
			fmt.Fprintln(out, 1)
			return
		}
		if a < n {
			less = append(less, a)
		} else if a > n {
			greater = append(greater, a)
		}
		if a < minA {
			minA = a
		}
		if a > maxA {
			maxA = a
		}
	}

	if n < minA || n > maxA || len(less) == 0 || len(greater) == 0 {
		fmt.Fprintln(out, -1)
		return
	}

	ans := int(^uint(0) >> 1)
	for _, ai := range less {
		for _, aj := range greater {
			g := gcd(n-ai, aj-n)
			L := (aj - ai) / g
			if L < ans {
				ans = L
			}
		}
	}

	fmt.Fprintln(out, ans)
}
