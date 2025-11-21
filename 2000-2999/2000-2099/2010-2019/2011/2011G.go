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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		ans := solveCase(p)
		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}

func solveCase(p []int) []int {
	n := len(p)
	pos := make([]int, n+1)
	for i, v := range p {
		pos[v] = i
	}

	prefInc := make([]int, n)
	prefDec := make([]int, n)
	for i := 0; i < n-1; i++ {
		prefInc[i+1] = prefInc[i]
		if p[i] < p[i+1] {
			prefInc[i+1]++
		}
		prefDec[i+1] = prefDec[i]
		if p[i] > p[i+1] {
			prefDec[i+1]++
		}
	}

	isInc := func(l, r int) bool {
		if l >= r {
			return true
		}
		return prefInc[r]-prefInc[l] == r-l
	}
	isDec := func(l, r int) bool {
		if l >= r {
			return true
		}
		return prefDec[r]-prefDec[l] == r-l
	}
	intervalCost := func(l, r int) int {
		if r-l <= 1 {
			return 0
		}
		L := l + 1
		R := r - 1
		if L > R {
			return 0
		}
		if isInc(L, R) || isDec(L, R) {
			return 1
		}
		return 2
	}

	ngl := make([]int, n)
	stack := make([]int, 0, n)
	for i := 0; i < n; i++ {
		for len(stack) > 0 && p[stack[len(stack)-1]] < p[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			ngl[i] = stack[len(stack)-1]
		} else {
			ngl[i] = -1
		}
		stack = append(stack, i)
	}

	ngr := make([]int, n)
	stack = stack[:0]
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && p[stack[len(stack)-1]] < p[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			ngr[i] = stack[len(stack)-1]
		} else {
			ngr[i] = -1
		}
		stack = append(stack, i)
	}

	res := make([]int, n-1)
	const inf = int(1e9)
	for val := 1; val < n; val++ {
		idx := pos[val]
		best := inf
		if left := ngl[idx]; left != -1 {
			if cost := intervalCost(left, idx) + 1; cost < best {
				best = cost
			}
		}
		if right := ngr[idx]; right != -1 {
			if cost := intervalCost(idx, right) + 1; cost < best {
				best = cost
			}
		}
		res[val-1] = best
	}
	return res
}
