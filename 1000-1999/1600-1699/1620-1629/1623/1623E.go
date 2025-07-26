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

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)

	left := make([]int, n+1)
	right := make([]int, n+1)
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &left[i], &right[i])
		if left[i] != 0 {
			parent[left[i]] = i
		}
		if right[i] != 0 {
			parent[right[i]] = i
		}
	}

	// in-order traversal
	order := make([]int, 0, n)
	stack := []int{}
	curr := 1
	for curr != 0 || len(stack) > 0 {
		for curr != 0 {
			stack = append(stack, curr)
			curr = left[curr]
		}
		curr = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, curr)
		curr = right[curr]
	}

	chars := []byte(s)

	// determine which positions are beneficial to duplicate
	want := make([]bool, n)
	i := 0
	for i < n {
		j := i + 1
		for j < n && chars[order[j]-1] == chars[order[i]-1] {
			j++
		}
		next := byte('{')
		if j < n {
			next = chars[order[j]-1]
		}
		if chars[order[i]-1] < next {
			for t := i; t < j; t++ {
				want[t] = true
			}
		}
		i = j
	}

	dup := make([]bool, n+1)
	dsu := make([]int, n+1)
	for i := 0; i <= n; i++ {
		dsu[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if dsu[x] != x {
			dsu[x] = find(dsu[x])
		}
		return dsu[x]
	}

	remaining := k
	for idx := 0; idx < n && remaining > 0; idx++ {
		u := order[idx]
		if dup[u] {
			continue
		}
		if want[idx] {
			path := make([]int, 0)
			v := find(u)
			for v != 0 && !dup[v] {
				path = append(path, v)
				v = find(parent[v])
			}
			if len(path) <= remaining {
				for _, x := range path {
					dup[x] = true
					dsu[x] = find(parent[x])
				}
				remaining -= len(path)
			}
		}
	}

	// build result string
	res := make([]byte, 0, n*2)
	stack = stack[:0]
	curr = 1
	for curr != 0 || len(stack) > 0 {
		for curr != 0 {
			stack = append(stack, curr)
			curr = left[curr]
		}
		curr = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		ch := chars[curr-1]
		res = append(res, ch)
		if dup[curr] {
			res = append(res, ch)
		}
		curr = right[curr]
	}

	fmt.Fprintln(out, string(res))
}
