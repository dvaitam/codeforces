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
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	rows := make(map[int][]int)
	cols := make(map[int][]int)

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			var c int
			fmt.Fscan(reader, &c)
			rows[c] = append(rows[c], i)
			cols[c] = append(cols[c], j)
		}
	}

	var ans int64
	for color, rSlice := range rows {
		sort.Ints(rSlice)
		prefix := 0
		for idx, v := range rSlice {
			ans += int64(v*idx - prefix)
			prefix += v
		}
		cSlice := cols[color]
		sort.Ints(cSlice)
		prefix = 0
		for idx, v := range cSlice {
			ans += int64(v*idx - prefix)
			prefix += v
		}
	}

	fmt.Fprintln(writer, ans)
}
