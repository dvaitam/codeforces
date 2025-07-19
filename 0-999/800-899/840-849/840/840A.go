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

	var m int
	_, _ = fmt.Fscan(reader, &m)
	a := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &a[i])
	}
	// b holds values and original indices
	b := make([]struct{ val, idx int }, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &b[i].val)
		b[i].idx = i
	}
	// sort a ascending and b by val ascending
	sort.Ints(a)
	sort.Slice(b, func(i, j int) bool {
		return b[i].val < b[j].val
	})
	// assign largest a to smallest b
	ans := make([]int, m)
	for i := 0; i < m; i++ {
		ans[b[i].idx] = a[m-i-1]
	}
	// output answer
	for i, v := range ans {
		if i > 0 {
			writer.WriteByte(' ')
		}
		writer.WriteString(fmt.Sprint(v))
	}
	writer.WriteByte('\n')
}
