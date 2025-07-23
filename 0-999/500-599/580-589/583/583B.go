package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	visited := make([]bool, n)
	pieces := 0
	count := 0
	pos := 0
	dir := 1
	changes := 0

	for count < n {
		if !visited[pos] && pieces >= a[pos] {
			visited[pos] = true
			pieces++
			count++
			if count == n {
				break
			}
		}
		next := pos + dir
		if next < 0 || next >= n {
			dir = -dir
			changes++
			next = pos + dir
		}
		pos = next
	}

	fmt.Fprintln(out, changes)
}
