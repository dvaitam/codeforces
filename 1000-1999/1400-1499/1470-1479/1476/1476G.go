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

	var n, m int
	fmt.Fscan(reader, &n, &m)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	for q := 0; q < m; q++ {
		var typ int
		fmt.Fscan(reader, &typ)
		if typ == 1 {
			var l, r, k int
			fmt.Fscan(reader, &l, &r, &k)
			l--
			r--
			cnt := make(map[int]int)
			for i := l; i <= r; i++ {
				cnt[a[i]]++
			}
			if len(cnt) < k {
				fmt.Fprintln(writer, -1)
				continue
			}
			freq := make([]int, 0, len(cnt))
			for _, v := range cnt {
				freq = append(freq, v)
			}
			sort.Ints(freq)
			best := int(1e9)
			for i := 0; i+k-1 < len(freq); i++ {
				diff := freq[i+k-1] - freq[i]
				if diff < best {
					best = diff
				}
			}
			fmt.Fprintln(writer, best)
		} else {
			var p, x int
			fmt.Fscan(reader, &p, &x)
			p--
			a[p] = x
		}
	}
}
