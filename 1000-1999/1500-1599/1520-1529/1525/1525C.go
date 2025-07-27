package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Robot struct {
	x   int
	dir byte
	idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		xs := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &xs[i])
		}
		dirs := make([]byte, n)
		var s string
		fmt.Fscan(in, &s)
		for i := 0; i < n; i++ {
			dirs[i] = s[i]
		}

		groups := [2][]Robot{}
		for i := 0; i < n; i++ {
			p := xs[i] % 2
			groups[p] = append(groups[p], Robot{xs[i], dirs[i], i})
		}

		ans := make([]int, n)
		for i := range ans {
			ans[i] = -1
		}

		for p := 0; p < 2; p++ {
			arr := groups[p]
			sort.Slice(arr, func(i, j int) bool { return arr[i].x < arr[j].x })
			stack := make([]Robot, 0)
			for _, r := range arr {
				if r.dir == 'R' {
					stack = append(stack, r)
				} else {
					if len(stack) > 0 && stack[len(stack)-1].dir == 'R' {
						prev := stack[len(stack)-1]
						stack = stack[:len(stack)-1]
						time := (r.x - prev.x) / 2
						ans[r.idx] = time
						ans[prev.idx] = time
					} else {
						r.x = -r.x
						r.dir = 'R'
						stack = append(stack, r)
					}
				}
			}
			for len(stack) >= 2 {
				a := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				b := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				time := (2*m - a.x - b.x) / 2
				ans[a.idx] = time
				ans[b.idx] = time
			}
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
