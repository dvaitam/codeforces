package main

import (
	"bufio"
	"fmt"
	"os"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign, val := 1, 0
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

func main() {
	fs := newFastScanner()
	n := fs.nextInt()
	m := fs.nextInt()
	k := fs.nextInt()

	products := make([]int, n+1)
	for i := 1; i <= n; i++ {
		products[i] = fs.nextInt()
	}

	graph := make([][]int, n+1)
	for i := 0; i < m; i++ {
		u := fs.nextInt()
		v := fs.nextInt()
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}

	dist := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = -1
	}

	queue := make([]int, 0, n)
	queue = append(queue, 1)
	dist[1] = 0

	for head := 0; head < len(queue); head++ {
		u := queue[head]
		for _, v := range graph[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				queue = append(queue, v)
			}
		}
	}

	ans := make([]int, k+1)
	for i := 1; i <= n; i++ {
		p := products[i]
		if dist[i] > ans[p] {
			ans[p] = dist[i]
		}
	}

	out := bufio.NewWriter(os.Stdout)
	for i := 1; i <= k; i++ {
		if i > 1 {
			fmt.Fprintf(out, " ")
		}
		fmt.Fprintf(out, "%d", ans[i])
	}
	fmt.Fprintln(out)
	out.Flush()
}
