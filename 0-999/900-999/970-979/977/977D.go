package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	arr := make([]int64, n)
	m := make(map[int64]int)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
		m[arr[i]] = i
	}
	indeg := make(map[int64]int)
	for _, v := range arr {
		if _, ok := m[v*2]; ok {
			indeg[v*2]++
		}
		if v%3 == 0 {
			if _, ok := m[v/3]; ok {
				indeg[v/3]++
			}
		}
	}
	var start int64
	for _, v := range arr {
		if indeg[v] == 0 {
			start = v
			break
		}
	}
	res := make([]int64, 0, n)
	used := make(map[int64]bool)
	var dfs func(int64)
	dfs = func(x int64) {
		res = append(res, x)
		used[x] = true
		if y, ok := m[x*2]; ok && !used[arr[y]] {
			dfs(arr[y])
		}
		if x%3 == 0 {
			if y, ok := m[x/3]; ok && !used[arr[y]] {
				dfs(arr[y])
			}
		}
	}
	dfs(start)
	writer := bufio.NewWriter(os.Stdout)
	for i, v := range res {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, v)
	}
	writer.WriteByte('\n')
	writer.Flush()
}
