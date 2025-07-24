package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		m := n / 2
		b := make([]int, m)
		used := make([]bool, n+1)
		ok := true
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &b[i])
			if b[i] < 1 || b[i] > n || used[b[i]] {
				ok = false
			}
			used[b[i]] = true
		}
		if !ok {
			fmt.Fprintln(out, -1)
			continue
		}
		avail := make([]int, 0, m)
		for i := 1; i <= n; i++ {
			if !used[i] {
				avail = append(avail, i)
			}
		}
		parent := make([]int, len(avail))
		for i := range parent {
			parent[i] = i
		}
		var find func(int) int
		find = func(x int) int {
			if x < 0 {
				return -1
			}
			if parent[x] == x {
				return x
			}
			parent[x] = find(parent[x])
			return parent[x]
		}
		remove := func(x int) { parent[x] = find(x - 1) }

		res := make([]int, n)
		for i := m - 1; i >= 0; i-- {
			idx := sort.SearchInts(avail, b[i]) - 1
			idx = find(idx)
			if idx < 0 {
				ok = false
				break
			}
			res[2*i] = avail[idx]
			res[2*i+1] = b[i]
			remove(idx)
		}
		if !ok {
			fmt.Fprintln(out, -1)
			continue
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, res[i])
		}
		fmt.Fprintln(out)
	}
}
