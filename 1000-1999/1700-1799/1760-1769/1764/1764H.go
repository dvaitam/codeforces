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

	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	L := make([]int, m)
	R := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &L[i], &R[i])
	}

	result := make([]int, m)
	for x := 0; x < m; x++ {
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = i + 1
		}
		next := make([]int, n+2)
		for i := 0; i <= n+1; i++ {
			next[i] = i
		}
		var find func(int) int
		find = func(p int) int {
			if p > n {
				return p
			}
			if next[p] != p {
				next[p] = find(next[p])
			}
			return next[p]
		}
		for t := 0; t < k; t++ {
			idx := (x + t) % m
			l := L[idx]
			r := R[idx]
			col := a[l-1]
			j := find(l + 1)
			for j <= r {
				a[j-1] = col
				next[j] = j + 1
				j = find(j + 1)
			}
		}
		cnt := 0
		prev := 0
		for i := 0; i < n; i++ {
			if i == 0 || a[i] != prev {
				cnt++
				prev = a[i]
			}
		}
		result[x] = cnt
	}

	for i := 0; i < m; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, result[i])
	}
	fmt.Fprintln(out)
}
