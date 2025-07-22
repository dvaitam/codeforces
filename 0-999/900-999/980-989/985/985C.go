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

	var n, k int
	var l int64
	if _, err := fmt.Fscan(reader, &n, &k, &l); err != nil {
		return
	}
	m := n * k
	a := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &a[i])
	}
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })

	limit := a[0] + l
	good := 0
	for good < m && a[good] <= limit {
		good++
	}
	if good < n {
		fmt.Fprintln(writer, 0)
		return
	}

	pos := good - 1
	right := m - 1
	var result int64
	for i := 0; i < n; i++ {
		result += a[pos]
		pos--
		for j := 0; j < k-1; j++ {
			if right > pos {
				right--
			} else {
				pos--
			}
		}
	}
	fmt.Fprintln(writer, result)
}
