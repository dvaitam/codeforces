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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		p := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &p[i])
		}

		d := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &d[i])
		}

		comp := make([]int, n+1)
		sizes := []int{0}
		visited := make([]bool, n+1)
		compID := 0

		for i := 1; i <= n; i++ {
			if visited[i] {
				continue
			}

			cur := []int{}
			x := i
			for !visited[x] {
				visited[x] = true
				cur = append(cur, x)
				x = p[x]
			}

			compID++
			sizes = append(sizes, len(cur))
			for _, v := range cur {
				comp[v] = compID
			}
		}

		active := make([]bool, compID+1)
		ans := make([]int64, n)
		var current int64

		for i := 0; i < n; i++ {
			idx := d[i+1]
			id := comp[idx]
			if !active[id] {
				active[id] = true
				current += int64(sizes[id])
			}
			ans[i] = current
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
