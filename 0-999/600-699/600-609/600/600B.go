package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	b := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &b[i])
	}
	sort.Ints(a)

	out := bufio.NewWriter(os.Stdout)
	for i := 0; i < m; i++ {
		if i > 0 {
			out.WriteByte(' ')
		}
		x := b[i]
		cnt := sort.Search(len(a), func(j int) bool { return a[j] > x })
		fmt.Fprint(out, cnt)
	}
	out.WriteByte('\n')
	out.Flush()
}
