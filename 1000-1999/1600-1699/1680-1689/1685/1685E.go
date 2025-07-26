package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func lis(a []int) int {
	d := make([]int, len(a))
	l := 0
	for _, x := range a {
		i := sort.Search(l, func(i int) bool { return d[i] >= x })
		if i == l {
			d[l] = x
			l++
		} else {
			d[i] = x
		}
	}
	return l
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	m := 2*n + 1
	p := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &p[i])
	}

	for ; q > 0; q-- {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		p[u], p[v] = p[v], p[u]

		ans := -1
		for k := 0; k < m; k++ {
			b := make([]int, m)
			for i := 0; i < m; i++ {
				b[i] = p[(k+i)%m]
			}
			if lis(b) <= n {
				ans = k
				break
			}
		}
		fmt.Fprintln(out, ans)
	}
}
