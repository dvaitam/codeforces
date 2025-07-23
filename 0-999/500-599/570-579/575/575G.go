package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
)

type Edge struct {
	to int
	d  int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)

	g := make([][]Edge, n)
	for i := 0; i < m; i++ {
		var a, b, l int
		fmt.Fscan(in, &a, &b, &l)
		g[a] = append(g[a], Edge{b, l})
		g[b] = append(g[b], Edge{a, l})
	}

	parent := make([]int, n)
	label := make([]int, n)
	vis := make([]bool, n)
	for i := range parent {
		parent[i] = -1
	}

	queue := [][]int{{n - 1}}
	vis[n-1] = true
	for len(queue) > 0 {
		curSet := queue[0]
		queue = queue[1:]
		buckets := make([][]int, 10)
		for _, v := range curSet {
			for _, e := range g[v] {
				if !vis[e.to] {
					vis[e.to] = true
					parent[e.to] = v
					label[e.to] = e.d
					buckets[e.d] = append(buckets[e.d], e.to)
				}
			}
		}
		for d := 0; d < 10; d++ {
			if len(buckets[d]) > 0 {
				queue = append(queue, buckets[d])
			}
		}
	}

	if !vis[0] {
		return
	}

	var digits []int
	var nodes []int
	cur := 0
	nodes = append(nodes, cur)
	for cur != n-1 {
		digits = append(digits, label[cur])
		cur = parent[cur]
		nodes = append(nodes, cur)
	}

	var sb strings.Builder
	for i := len(digits) - 1; i >= 0; i-- {
		sb.WriteByte(byte('0' + digits[i]))
	}
	tStr := sb.String()
	if tStr == "" {
		tStr = "0"
	}
	var bigT big.Int
	bigT.SetString(tStr, 10)

	fmt.Fprintln(out, bigT.String())
	fmt.Fprintln(out, len(nodes))
	for i, v := range nodes {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
