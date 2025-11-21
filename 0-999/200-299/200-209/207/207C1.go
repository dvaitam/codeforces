package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	ch byte
}

var (
	tree    [2][][]Edge
	nodeCnt [2]int
	freq1   map[string]int64
	freq2   map[string]int64
)

func dfs1(v int, rev string, t int) {
	freq1[rev]++
	for _, e := range tree[t][v] {
		dfs1(e.to, string(e.ch)+rev, t)
	}
}

func dfs2(v int) map[string]int64 {
	res := make(map[string]int64)
	res[""] = 1
	freq2[""]++
	for _, e := range tree[1][v] {
		childMap := dfs2(e.to)
		for s, cnt := range childMap {
			newStr := string(e.ch) + s
			res[newStr] += cnt
			freq2[newStr] += cnt
		}
	}
	return res
}

func computeAns() int64 {
	freq1 = make(map[string]int64)
	freq2 = make(map[string]int64)
	dfs1(1, "", 0)
	dfs2(1)
	var ans int64
	for s, c := range freq1 {
		if c2, ok := freq2[s]; ok {
			ans += c * c2
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	maxNodes := n + 5
	for t := 0; t < 2; t++ {
		tree[t] = make([][]Edge, maxNodes)
		nodeCnt[t] = 1
	}

	for i := 0; i < n; i++ {
		var t, v int
		var cs string
		fmt.Fscan(in, &t, &v, &cs)
		t--
		nodeCnt[t]++
		child := nodeCnt[t]
		tree[t][v] = append(tree[t][v], Edge{to: child, ch: cs[0]})
		fmt.Fprintln(out, computeAns())
	}
}
