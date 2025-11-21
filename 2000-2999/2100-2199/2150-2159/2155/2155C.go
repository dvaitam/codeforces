package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 676767677

func validate(b []int, n, totalLeft int) bool {
	if totalLeft < 0 || totalLeft > n {
		return false
	}
	prefix := 0
	for i := 0; i < n; i++ {
		left := b[i] + totalLeft - 2*prefix
		if left != 0 && left != 1 {
			return false
		}
		prefix += left
		if prefix < 0 || prefix > i+1 {
			return false
		}
	}
	return prefix == totalLeft
}

func contains(arr []int, v int) bool {
	for _, x := range arr {
		if x == v {
			return true
		}
	}
	return false
}

func solve(n int, a []int) int {
	b := make([]int, n)
	s := make([]int, n+1)
	for i := 1; i <= n; i++ {
		idx := i - 1
		b[idx] = a[idx] - (n - i + 1)
		s[i] = b[idx] - s[i-1]
	}

	if n%2 == 0 {
		if validate(b, n, s[n]) {
			return 1
		}
		return 0
	}

	if s[n] != 0 {
		return 0
	}

	candidates := []int{}
	for i := 1; i <= n; i++ {
		delta := s[i] - s[i-1]
		var opts []int
		if i%2 == 1 {
			x := -delta
			y := 1 - delta
			if x > y {
				x, y = y, x
			}
			opts = append(opts, x)
			if y != x {
				opts = append(opts, y)
			}
		} else {
			x := delta - 1
			y := delta
			if x > y {
				x, y = y, x
			}
			opts = append(opts, x)
			if y != x {
				opts = append(opts, y)
			}
		}

		if len(candidates) == 0 {
			candidates = append(candidates, opts...)
			continue
		}

		next := make([]int, 0, len(candidates))
		for _, v := range candidates {
			if contains(opts, v) && !contains(next, v) {
				next = append(next, v)
			}
		}
		candidates = next
		if len(candidates) == 0 {
			break
		}
	}

	ans := 0
	for _, T := range candidates {
		if validate(b, n, T) {
			ans++
			if ans >= mod {
				ans -= mod
			}
		}
	}
	return ans
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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		fmt.Fprintln(out, solve(n, a)%mod)
	}
}
