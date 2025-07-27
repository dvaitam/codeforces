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
	L := make([]int, m)
	R := make([]int, m)
	Len := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &L[i], &R[i])
		Len[i] = R[i] - L[i] + 1
	}

	for x := 1; x <= n; x++ {
		parent := make([]int, n+2)
		for i := 0; i <= n+1; i++ {
			parent[i] = i
		}
		var find func(int) int
		find = func(a int) int {
			if parent[a] != a {
				parent[a] = find(parent[a])
			}
			return parent[a]
		}
		occupy := func(l, r int) {
			i := find(l)
			for i <= r {
				parent[i] = find(i + 1)
				i = parent[i]
			}
		}
		total := 0
		for i := 0; i < m; i++ {
			if Len[i] < x {
				continue
			}
			if find(L[i]) > R[i] {
				occupy(L[i], R[i])
				total += Len[i]
			}
		}
		fmt.Fprintln(writer, total)
	}
}
