package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	var k int64
	fmt.Fscan(in, &n, &k)
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	if k == 1 {
		// pick minimal element
		idx := 0
		for i := 1; i < n; i++ {
			if a[i] < a[idx] {
				idx = i
			}
		}
		fmt.Fprintln(out, 1)
		// 1-based
		fmt.Fprintln(out, idx+1)
		return
	}
	// factorize k into primes and exponents
	pr := []int64{}
	cnt := []int{}
	kk := k
	for p := int64(2); p*p <= kk; p++ {
		if kk%p == 0 {
			e := 0
			for kk%p == 0 {
				kk /= p
				e++
			}
			pr = append(pr, p)
			cnt = append(cnt, e)
		}
	}
	if kk > 1 {
		pr = append(pr, kk)
		cnt = append(cnt, 1)
	}
	m := len(pr)
	// compute total states
	maxstate := 1
	for j := 0; j < m; j++ {
		maxstate *= (cnt[j] + 1)
	}
	// DP arrays
	inf := n + 1
	f := make([]int, maxstate)
	g := make([]int64, maxstate)
	represent := make([][]int, maxstate)
	List := make([][]int, maxstate)
	for i := 0; i < maxstate; i++ {
		f[i] = inf
	}
	// initial state: full exponents needed
	start := maxstate - 1
	f[start] = 0
	g[start] = 0
	// represent[start] = cnt copy
	represent[start] = make([]int, m)
	for j := 0; j < m; j++ {
		represent[start][j] = cnt[j]
	}
	// DP
	for i := 0; i < n; i++ {
		// compute current exponents in a[i]
		x := a[i]
		currExp := make([]int, m)
		for j := 0; j < m; j++ {
			for x%pr[j] == 0 {
				x /= pr[j]
				currExp[j]++
			}
		}
		// try transitions
		for curr := 0; curr < maxstate; curr++ {
			if f[curr] == inf {
				continue
			}
			// compute next state
			next := 0
			// need represent[curr]
			repc := represent[curr]
			nextRepr := make([]int, m)
			for j := 0; j < m; j++ {
				rem := repc[j] - currExp[j]
				if rem < 0 {
					rem = 0
				}
				// build mixed-radix to index
				nextRepr[j] = rem
			}
			// convert nextRepr to index
			mul := 1
			for j := 0; j < m; j++ {
				next += nextRepr[j] * mul
				mul *= (cnt[j] + 1)
			}
			// relax
			candF := f[curr] + 1
			candG := g[curr] + a[i]
			if candF < f[next] || (candF == f[next] && candG < g[next]) {
				f[next] = candF
				g[next] = candG
				// update represent and List
				represent[next] = nextRepr
				// copy List[curr] and append i
				lst := List[curr]
				newList := make([]int, len(lst)+1)
				copy(newList, lst)
				newList[len(lst)] = i
				List[next] = newList
			}
		}
	}
	// output
	if f[0] == inf {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, f[0])
		// print 1-based indices
		for _, idx := range List[0] {
			fmt.Fprint(out, idx+1, " ")
		}
		fmt.Fprintln(out)
	}
}
