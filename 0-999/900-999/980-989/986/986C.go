package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n     int
	exist []bool
	vis   []bool
	seen  []bool
	full  int
)

func dfsNode(v int) {
	if vis[v] {
		return
	}
	vis[v] = true
	dfsMask(full ^ v)
}

func dfsMask(mask int) {
	if seen[mask] {
		return
	}
	seen[mask] = true
	if exist[mask] && !vis[mask] {
		dfsNode(mask)
	}
	for i := 0; i < n; i++ {
		if mask&(1<<uint(i)) != 0 {
			dfsMask(mask &^ (1 << uint(i)))
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	size := 1 << uint(n)
	exist = make([]bool, size)
	vis = make([]bool, size)
	seen = make([]bool, size)
	full = size - 1

	nums := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &nums[i])
		exist[nums[i]] = true
	}

	components := 0
	for _, v := range nums {
		if !vis[v] {
			components++
			dfsNode(v)
		}
	}

	fmt.Fprintln(writer, components)
}
