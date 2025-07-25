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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	pos := make([]int, n+1)
	posts := make([]int, n+1)
	minPos := make([]int, n+1)
	maxPos := make([]int, n+1)

	for i := 1; i <= n; i++ {
		pos[i] = i
		posts[i] = i
		minPos[i] = i
		maxPos[i] = i
	}

	for j := 0; j < m; j++ {
		var x int
		fmt.Fscan(in, &x)
		if pos[x] > 1 {
			y := posts[pos[x]-1]
			posts[pos[x]-1], posts[pos[x]] = posts[pos[x]], posts[pos[x]-1]
			pos[x]--
			pos[y]++

			if minPos[x] > pos[x] {
				minPos[x] = pos[x]
			}
			if maxPos[y] < pos[y] {
				maxPos[y] = pos[y]
			}
		}
		if maxPos[x] < pos[x] {
			maxPos[x] = pos[x]
		}
	}

	for i := 1; i <= n; i++ {
		fmt.Fprintln(out, minPos[i], maxPos[i])
	}
}
