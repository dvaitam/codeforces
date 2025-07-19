package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for T > 0 {
		T--
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n+1)
		c := make([]int, n+1)
		in := make([]int, n+1)
		f := make([]bool, n+1)
		v := make([][]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
			v[i] = append(v[i], a[i])
			in[a[i]]++
		}
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &c[i])
		}
		// topological removal of non-cycle nodes
		ans := make([]int, 0, n)
		queue := make([]int, 0, n)
		head := 0
		for i := 1; i <= n; i++ {
			if in[i] == 0 {
				queue = append(queue, i)
				ans = append(ans, i)
				f[i] = true
			}
		}
		for head < len(queue) {
			t := queue[head]
			head++
			for _, to := range v[t] {
				in[to]--
				if in[to] == 0 {
					queue = append(queue, to)
					ans = append(ans, to)
					f[to] = true
				}
			}
		}
		// collect remaining cycle nodes
		others := make([]int, 0, n)
		for i := 1; i <= n; i++ {
			if !f[i] {
				others = append(others, i)
			}
		}
		sort.Slice(others, func(i, j int) bool {
			return c[others[i]] < c[others[j]]
		})
		// process each cycle
		for _, t := range others {
			if f[t] {
				continue
			}
			// mark start and traverse cycle
			f[t] = true
			cur := a[t]
			for cur != t {
				ans = append(ans, cur)
				f[cur] = true
				cur = a[cur]
			}
			ans = append(ans, t)
		}
		// output answer
		for i, x := range ans {
			if i > 0 {
				writer.WriteByte(' ')
			}
			writer.WriteString(fmt.Sprintf("%d", x))
		}
		writer.WriteByte('\n')
	}
}
