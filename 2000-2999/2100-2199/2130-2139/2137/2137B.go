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
	if a < 0 {
		return -a
	}
	return a
}

func assignSimple(p []int, base [][]int) []int {
	n := len(p)
	lens := []int{len(base[0]), len(base[1]), len(base[2])}
	q := make([]int, n)
	for i := 0; i < n; i++ {
		target := (3 - (p[i] % 3)) % 3
		if lens[target] == 0 {
			panic("insufficient numbers")
		}
		lens[target]--
		q[i] = base[target][lens[target]]
	}
	return q
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}

		base := make([][]int, 3)
		for x := 1; x <= n; x++ {
			r := x % 3
			base[r] = append(base[r], x)
		}

		if n%3 != 1 {
			q := assignSimple(p, base)
			for i, val := range q {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, val)
			}
			fmt.Fprintln(out)
			continue
		}

		patterns := make([][4]int, 0, 81)
		for a := 0; a < 3; a++ {
			for b := 0; b < 3; b++ {
				for c := 0; c < 3; c++ {
					for d := 0; d < 3; d++ {
						patterns = append(patterns, [4]int{a, b, c, d})
					}
				}
			}
		}

		success := false
		for _, pattern := range patterns {
			lens := []int{len(base[0]), len(base[1]), len(base[2])}
			pool := make([]int, 0, 4)
			ok := true
			for _, r := range pattern {
				if lens[r] == 0 {
					ok = false
					break
				}
				lens[r]--
				pool = append(pool, base[r][lens[r]])
			}
			if !ok {
				continue
			}

			q := make([]int, n)
			for i := 4; i < n; i++ {
				target := (3 - (p[i] % 3)) % 3
				if lens[target] == 0 {
					ok = false
					break
				}
				lens[target]--
				q[i] = base[target][lens[target]]
			}
			if !ok {
				continue
			}

			used := make([]bool, 4)
			idx := []int{0, 1, 2, 3}
			found := false
			var dfs func(int)
			dfs = func(pos int) {
				if found {
					return
				}
				if pos == len(idx) {
					for i := 0; i < n-1; i++ {
						if gcd(p[i]+q[i], p[i+1]+q[i+1]) < 3 {
							return
						}
					}
					found = true
					return
				}
				cur := idx[pos]
				for j := 0; j < len(pool); j++ {
					if used[j] {
						continue
					}
					q[cur] = pool[j]
					used[j] = true
					dfs(pos + 1)
					if found {
						return
					}
					used[j] = false
				}
			}
			dfs(0)
			if found {
				for i, val := range q {
					if i > 0 {
						fmt.Fprint(out, " ")
					}
					fmt.Fprint(out, val)
				}
				fmt.Fprintln(out)
				success = true
				break
			}
		}
		if !success {
			panic("no valid pattern")
		}
	}
}
