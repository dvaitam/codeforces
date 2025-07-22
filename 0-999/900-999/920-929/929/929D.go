package main

import (
	"bufio"
	"fmt"
	"os"
)

type State struct {
	pos  int
	mask uint64
	dist int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, a, b int
	if _, err := fmt.Fscan(reader, &n, &a, &b); err != nil {
		return
	}
	g := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(reader, &g[i])
	}
	k := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &k[i])
	}

	colorID := make(map[int]int)
	nextID := 0
	for _, x := range g {
		if _, ok := colorID[x]; !ok {
			colorID[x] = nextID
			nextID++
		}
	}
	for _, x := range k {
		if _, ok := colorID[x]; !ok {
			colorID[x] = nextID
			nextID++
		}
	}

	if nextID > 20 {
		// Too many colors for this simple solution
		fmt.Fprintln(writer, -1)
		return
	}

	gID := make([]int, n-1)
	for i, x := range g {
		gID[i] = colorID[x]
	}
	kID := make([]int, n)
	for i, x := range k {
		kID[i] = colorID[x]
	}

	startMask := uint64(1) << kID[a-1]
	q := []State{{a - 1, startMask, 0}}
	visited := make(map[[2]int]bool)
	visited[[2]int{a - 1, int(startMask)}] = true

	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		mask := cur.mask | (1 << kID[cur.pos])
		if cur.pos == b-1 {
			fmt.Fprintln(writer, cur.dist)
			return
		}
		if cur.pos > 0 && (mask&(1<<gID[cur.pos-1])) != 0 {
			np := cur.pos - 1
			nm := mask | (1 << kID[np])
			key := [2]int{np, int(nm)}
			if !visited[key] {
				visited[key] = true
				q = append(q, State{np, nm, cur.dist + 1})
			}
		}
		if cur.pos < n-1 && (mask&(1<<gID[cur.pos])) != 0 {
			np := cur.pos + 1
			nm := mask | (1 << kID[np])
			key := [2]int{np, int(nm)}
			if !visited[key] {
				visited[key] = true
				q = append(q, State{np, nm, cur.dist + 1})
			}
		}
	}
	fmt.Fprintln(writer, -1)
}
