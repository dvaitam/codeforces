package main

import (
	"bufio"
	"fmt"
	"os"
)

var parent []int

func find(x int) int {
	for x > 0 && parent[x] != x {
		parent[x] = parent[parent[x]]
		x = parent[x]
	}
	if x < 0 {
		return 0
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	pref := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &pref[i])
	}

	parent = make([]int, n+1)
	for i := 0; i <= n; i++ {
		parent[i] = i
	}

	stack := make([]int, 0)
	need := 1
	last := n + 1
	res := make([]int, n)

	for i := 0; i < k; i++ {
		x := pref[i]
		if x >= last {
			fmt.Fprintln(out, -1)
			return
		}
		res[i] = x
		parent[x] = find(x - 1)
		stack = append(stack, x)
		last = x
		for len(stack) > 0 && stack[len(stack)-1] == need {
			stack = stack[:len(stack)-1]
			need++
			if len(stack) > 0 {
				last = stack[len(stack)-1]
			} else {
				last = n + 1
			}
		}
	}

	for i := k; i < n; i++ {
		for len(stack) > 0 && stack[len(stack)-1] == need {
			stack = stack[:len(stack)-1]
			need++
			if len(stack) > 0 {
				last = stack[len(stack)-1]
			} else {
				last = n + 1
			}
		}
		x := find(last - 1)
		if x == 0 {
			fmt.Fprintln(out, -1)
			return
		}
		res[i] = x
		parent[x] = find(x - 1)
		stack = append(stack, x)
		last = x
	}

	for i, v := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
