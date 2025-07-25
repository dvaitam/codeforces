package main

import (
	"bufio"
	"fmt"
	"os"
)

const ALPHA = 26

type Node struct {
	next [ALPHA]int
	fail int
	out  []int
}

var nodes []Node

func addString(s string, id int) {
	v := 0
	for i := 0; i < len(s); i++ {
		c := int(s[i] - 'A')
		if nodes[v].next[c] == 0 {
			nodes = append(nodes, Node{})
			nodes[v].next[c] = len(nodes) - 1
		}
		v = nodes[v].next[c]
	}
	nodes[v].out = append(nodes[v].out, id)
}

func buildFail() {
	queue := make([]int, 0)
	for c := 0; c < ALPHA; c++ {
		if nodes[0].next[c] != 0 {
			queue = append(queue, nodes[0].next[c])
		}
	}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		f := nodes[v].fail
		for c := 0; c < ALPHA; c++ {
			u := nodes[v].next[c]
			if u != 0 {
				nodes[u].fail = nodes[f].next[c]
				queue = append(queue, u)
			}
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var R, C, Q int
	if _, err := fmt.Fscan(in, &R, &C, &Q); err != nil {
		return
	}
	grid := make([][]byte, R)
	for i := 0; i < R; i++ {
		var line string
		fmt.Fscan(in, &line)
		grid[i] = []byte(line)
	}
	queries := make([]string, Q)
	for i := 0; i < Q; i++ {
		fmt.Fscan(in, &queries[i])
	}
	nodes = make([]Node, 1)
	for i, s := range queries {
		addString(s, i)
	}
	buildFail()
	res := make([]int, Q)
	vertical := make([]map[int]int, C)
	for j := 0; j < C; j++ {
		vertical[j] = make(map[int]int)
	}
	for r := 0; r < R; r++ {
		// extend vertical
		for j := 0; j < C; j++ {
			ch := grid[r][j] - 'A'
			newMap := make(map[int]int)
			for st, cnt := range vertical[j] {
				ns := nodes[st].next[ch]
				if ns != 0 {
					newMap[ns] += cnt
				}
			}
			for st, cnt := range newMap {
				for _, idx := range nodes[st].out {
					res[idx] += cnt
				}
			}
			vertical[j] = newMap
		}
		// horizontal scanning
		active := make(map[int]int)
		for j := 0; j < C; j++ {
			ch := grid[r][j] - 'A'
			activeNext := make(map[int]int)
			for st, cnt := range active {
				ns := nodes[st].next[ch]
				if ns != 0 {
					activeNext[ns] += cnt
				}
			}
			ns := nodes[0].next[ch]
			if ns != 0 {
				activeNext[ns]++
			}
			for st, cnt := range activeNext {
				for _, idx := range nodes[st].out {
					res[idx] += cnt
				}
				vertical[j][st] += cnt
			}
			active = activeNext
		}
	}
	for i := 0; i < Q; i++ {
		fmt.Fprintln(out, res[i])
	}
}
