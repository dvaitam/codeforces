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
	sign := 1
	val := 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func main() {
	fs := newFastScanner()
	n := fs.nextInt()
	m := fs.nextInt()

	adj := make([][]int, n)
	radj := make([][]int, n)
	for i := 0; i < m; i++ {
		u := fs.nextInt() - 1
		v := fs.nextInt() - 1
		adj[u] = append(adj[u], v)
		radj[v] = append(radj[v], u)
	}

	k := fs.nextInt()
	path := make([]int, k)
	for i := 0; i < k; i++ {
		path[i] = fs.nextInt() - 1
	}

	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	queue := make([]int, 0, n)
	t := path[k-1]
	dist[t] = 0
	queue = append(queue, t)
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		for _, to := range radj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				queue = append(queue, to)
			}
		}
	}

	countShort := make([]int, n)
	for u := 0; u < n; u++ {
		if dist[u] == -1 {
			continue
		}
		cnt := 0
		for _, v := range adj[u] {
			if dist[v] == dist[u]-1 {
				cnt++
			}
		}
		countShort[u] = cnt
	}

	minRebuild := 0
	maxRebuild := 0
	for i := 0; i < k-1; i++ {
		u := path[i]
		v := path[i+1]
		if dist[u] != dist[v]+1 {
			minRebuild++
			maxRebuild++
		} else if countShort[u] >= 2 {
			maxRebuild++
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintf(writer, "%d %d\n", minRebuild, maxRebuild)
	writer.Flush()
}
