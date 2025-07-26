package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	start    int
	end      int
	parent   int
	childIdx int
	children []int
	pref     []int64
	count    int64
}

func buildTree(s string) ([]Node, []int, []int, []int) {
	n := len(s)
	matchOpen := make([]int, n+1)
	matchClose := make([]int, n+1)
	stack := make([]int, 0)
	for i := 1; i <= n; i++ {
		if s[i-1] == '(' {
			stack = append(stack, i)
		} else {
			if len(stack) > 0 {
				open := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				matchOpen[open] = i
				matchClose[i] = open
			}
		}
	}

	nodes := make([]Node, 1, n/2+1) // node 0 is root
	startToID := make([]int, n+1)
	nodeStack := make([]int, 0)
	for i := 1; i <= n; i++ {
		if s[i-1] == '(' {
			if matchOpen[i] == 0 {
				continue
			}
			id := len(nodes)
			nodes = append(nodes, Node{start: i})
			startToID[i] = id
			nodeStack = append(nodeStack, id)
		} else {
			openIdx := matchClose[i]
			if openIdx == 0 {
				continue
			}
			if len(nodeStack) == 0 {
				continue
			}
			id := nodeStack[len(nodeStack)-1]
			nodeStack = nodeStack[:len(nodeStack)-1]
			nodes[id].end = i
			parent := 0
			if len(nodeStack) > 0 {
				parent = nodeStack[len(nodeStack)-1]
			}
			nodes[id].parent = parent
			nodes[id].childIdx = len(nodes[parent].children)
			nodes[parent].children = append(nodes[parent].children, id)
			k := len(nodes[id].children)
			cnt := int64(1)
			for _, ch := range nodes[id].children {
				cnt += nodes[ch].count
			}
			cnt += int64(k * (k - 1) / 2)
			nodes[id].count = cnt
		}
	}
	for i := range nodes {
		pref := make([]int64, len(nodes[i].children)+1)
		for j, ch := range nodes[i].children {
			pref[j+1] = pref[j] + nodes[ch].count
		}
		nodes[i].pref = pref
	}
	return nodes, matchOpen, matchClose, startToID
}

func query(nodes []Node, rev []int, startToID []int, l, r int) int64 {
	openR := rev[r]
	nodeL := startToID[l]
	nodeR := startToID[openR]
	parent := nodes[nodeL].parent
	idxL := nodes[nodeL].childIdx
	idxR := nodes[nodeR].childIdx
	pref := nodes[parent].pref
	t := idxR - idxL + 1
	sumCounts := pref[idxR+1] - pref[idxL]
	return sumCounts + int64(t*(t-1)/2)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	var s string
	fmt.Fscan(in, &s)

	nodes, _, rev, startToID := buildTree(s)

	for i := 0; i < q; i++ {
		var t, l, r int
		fmt.Fscan(in, &t, &l, &r)
		ans := query(nodes, rev, startToID, l, r)
		fmt.Fprintln(out, ans)
	}
}
