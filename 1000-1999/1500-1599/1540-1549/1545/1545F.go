package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	a := make([]int, n+1)
	b := make([]int, n+1)
	c := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &c[i])
	}

	// temporary counts reused for each query
	prefixB := make([]int, n+1)
	suffixC := make([]int, n+1)
	prefixUsed := make([]int, 0)
	suffixUsed := make([]int, 0)

	for op := 0; op < m; op++ {
		var typ int
		fmt.Fscan(reader, &typ)
		if typ == 1 {
			var k, x int
			fmt.Fscan(reader, &k, &x)
			a[k] = x
		} else if typ == 2 {
			var r int
			fmt.Fscan(reader, &r)
			prefixUsed = prefixUsed[:0]
			suffixUsed = suffixUsed[:0]
			// build suffix counts of c[a[k]] for k=1..r
			for i := r; i >= 1; i-- {
				val := c[a[i]]
				if suffixC[val] == 0 {
					suffixUsed = append(suffixUsed, val)
				}
				suffixC[val]++
			}
			var ans int64
			for j := 1; j <= r; j++ {
				val := c[a[j]]
				suffixC[val]--
				x := a[j]
				ans += int64(prefixB[x]) * int64(suffixC[x])
				bv := b[a[j]]
				if prefixB[bv] == 0 {
					prefixUsed = append(prefixUsed, bv)
				}
				prefixB[bv]++
			}
			// reset used counters
			for _, v := range prefixUsed {
				prefixB[v] = 0
			}
			for _, v := range suffixUsed {
				suffixC[v] = 0
			}
			fmt.Fprintln(writer, ans)
		}
	}
}
