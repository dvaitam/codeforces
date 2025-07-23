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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	dist := make([]int, n+1)
	for i := 2; i <= n; i++ {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	q = append(q, 1)
	head := 0
	for head < len(q) {
		v := q[head]
		head++
		d := dist[v]
		// move to v-1
		if v > 1 && dist[v-1] == -1 {
			dist[v-1] = d + 1
			q = append(q, v-1)
		}
		// move to v+1
		if v < n && dist[v+1] == -1 {
			dist[v+1] = d + 1
			q = append(q, v+1)
		}
		// shortcut to a[v]
		if dist[a[v]] == -1 {
			dist[a[v]] = d + 1
			q = append(q, a[v])
		}
	}

	for i := 1; i <= n; i++ {
		if i > 1 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, dist[i])
	}
}
