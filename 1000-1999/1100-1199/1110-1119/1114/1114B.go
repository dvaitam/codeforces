package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Node struct {
	val int64
	pos int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, k int
	fmt.Fscan(reader, &n, &m, &k)
	a := make([]Node, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i].val)
		a[i].pos = i
	}
	b := make([]Node, n)
	copy(b, a)

	// get top m*k elements
	tar := m * k
	sort.Slice(b, func(i, j int) bool {
		return b[i].val > b[j].val
	})
	vis := make([]bool, n)
	var sum int64
	for i := 0; i < tar; i++ {
		sum += b[i].val
		vis[b[i].pos] = true
	}
	// output sum
	fmt.Fprintln(writer, sum)

	// output k-1 cut positions
	cnt, t := 0, 0
	for i := 0; i < n; i++ {
		if t == k-1 {
			break
		}
		if vis[i] {
			cnt++
		}
		if cnt == m {
			// print position i+1
			fmt.Fprint(writer, i+1)
			if t < k-2 {
				fmt.Fprint(writer, " ")
			}
			cnt = 0
			t++
		}
	}
	if k > 1 {
		fmt.Fprintln(writer)
	}
}
