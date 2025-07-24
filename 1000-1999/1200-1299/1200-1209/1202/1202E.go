package main

import (
	"bufio"
	"fmt"
	"os"
)

type node struct {
	next [26]int
	fail int
	out  int
}

func newNode() node {
	n := node{fail: 0, out: 0}
	for i := 0; i < 26; i++ {
		n.next[i] = -1
	}
	return n
}

func buildAC(patterns []string) []node {
	nodes := make([]node, 1)
	nodes[0] = newNode()
	for _, s := range patterns {
		v := 0
		for i := 0; i < len(s); i++ {
			c := int(s[i] - 'a')
			if nodes[v].next[c] == -1 {
				nodes = append(nodes, newNode())
				nodes[v].next[c] = len(nodes) - 1
			}
			v = nodes[v].next[c]
		}
		nodes[v].out++
	}
	q := make([]int, 0)
	for c := 0; c < 26; c++ {
		v := nodes[0].next[c]
		if v != -1 {
			nodes[v].fail = 0
			q = append(q, v)
		} else {
			nodes[0].next[c] = 0
		}
	}
	for idx := 0; idx < len(q); idx++ {
		v := q[idx]
		f := nodes[v].fail
		nodes[v].out += nodes[f].out
		for c := 0; c < 26; c++ {
			u := nodes[v].next[c]
			if u != -1 {
				nodes[u].fail = nodes[f].next[c]
				q = append(q, u)
			} else {
				nodes[v].next[c] = nodes[f].next[c]
			}
		}
	}
	return nodes
}

func countEnds(nodes []node, text string) []int {
	n := len(text)
	res := make([]int, n)
	state := 0
	for i := 0; i < n; i++ {
		c := int(text[i] - 'a')
		state = nodes[state].next[c]
		res[i] = nodes[state].out
	}
	return res
}

func countStarts(nodes []node, text string) []int {
	n := len(text)
	res := make([]int, n)
	state := 0
	for i := n - 1; i >= 0; i-- {
		c := int(text[i] - 'a')
		state = nodes[state].next[c]
		res[i] = nodes[state].out
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t string
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	var n int
	fmt.Fscan(in, &n)
	strs := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &strs[i])
	}

	acForward := buildAC(strs)
	rev := make([]string, n)
	for i := 0; i < n; i++ {
		b := []byte(strs[i])
		for l, r := 0, len(b)-1; l < r; l, r = l+1, r-1 {
			b[l], b[r] = b[r], b[l]
		}
		rev[i] = string(b)
	}
	acBackward := buildAC(rev)

	pref := countEnds(acForward, t)
	suff := countStarts(acBackward, t)

	var ans int64
	for i := 0; i < len(t)-1; i++ {
		ans += int64(pref[i]) * int64(suff[i+1])
	}
	fmt.Fprintln(out, ans)
}
