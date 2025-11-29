package main

import (
	"bufio"
	"os"
	"strconv"
)

const MAX_BITS = 17

type Node struct {
	l, r, cnt int
}

var nodes []Node
var nodeIdx int

func newNode() int {
	if nodeIdx >= len(nodes) {
		nodes = append(nodes, Node{l: -1, r: -1, cnt: 0})
	} else {
		nodes[nodeIdx] = Node{l: -1, r: -1, cnt: 0}
	}
	idx := nodeIdx
	nodeIdx++
	return idx
}

func insert(root int, val int) {
	curr := root
	nodes[curr].cnt++
	for i := MAX_BITS - 1; i >= 0; i-- {
		bit := (val >> i) & 1
		if bit == 0 {
			if nodes[curr].l == -1 {
				nodes[curr].l = newNode()
			}
			curr = nodes[curr].l
		} else {
			if nodes[curr].r == -1 {
				nodes[curr].r = newNode()
			}
			curr = nodes[curr].r
		}
		nodes[curr].cnt++
	}
}

func getCnt(idx int) int {
	if idx == -1 {
		return 0
	}
	return nodes[idx].cnt
}

type Pair struct {
	u, v int
}

func solve(bit int, pairs []Pair) (int, bool) {
	if bit < 0 {
		return 0, true
	}

	// Try bit = 0
	possible0 := true
	var nextPairs0 []Pair
	nextPairs0 = make([]Pair, 0, len(pairs))

	for _, p := range pairs {
		u, v := p.u, p.v
		cntUL, cntUR := getCnt(nodes[u].l), getCnt(nodes[u].r)
		cntVL, cntVR := getCnt(nodes[v].l), getCnt(nodes[v].r)

		if cntUL != cntVL || cntUR != cntVR {
			possible0 = false
			break
		}
		if nodes[u].l != -1 {
			nextPairs0 = append(nextPairs0, Pair{nodes[u].l, nodes[v].l})
		}
		if nodes[u].r != -1 {
			nextPairs0 = append(nextPairs0, Pair{nodes[u].r, nodes[v].r})
		}
	}

	if possible0 {
		res, ok := solve(bit-1, nextPairs0)
		if ok {
			return res, true
		}
	}

	// Try bit = 1
	possible1 := true
	var nextPairs1 []Pair
	nextPairs1 = make([]Pair, 0, len(pairs))

	for _, p := range pairs {
		u, v := p.u, p.v
		cntUL, cntUR := getCnt(nodes[u].l), getCnt(nodes[u].r)
		cntVL, cntVR := getCnt(nodes[v].l), getCnt(nodes[v].r)

		if cntUL != cntVR || cntUR != cntVL {
			possible1 = false
			break
		}
		if nodes[u].l != -1 {
			nextPairs1 = append(nextPairs1, Pair{nodes[u].l, nodes[v].r})
		}
		if nodes[u].r != -1 {
			nextPairs1 = append(nextPairs1, Pair{nodes[u].r, nodes[v].l})
		}
	}

	if possible1 {
		res, ok := solve(bit-1, nextPairs1)
		if ok {
			return res | (1 << bit), true
		}
	}

	return 0, false
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 1024*1024)
	scanner.Split(bufio.ScanWords)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	// Preallocate a large buffer for nodes to avoid overhead and GC
	// Max nodes approx 17 bits * 131072 * 2 per test case, reused.
	// 5 million is safe.
	nodes = make([]Node, 5000000)

	scanner.Scan()
	tStr := scanner.Text()
	if tStr == "" {
		return
	}
	t, _ := strconv.Atoi(tStr)

	for i := 0; i < t; i++ {
		scanner.Scan()
		l, _ := strconv.Atoi(scanner.Text())
		scanner.Scan()
		r, _ := strconv.Atoi(scanner.Text())

		nodeIdx = 0
		rootS := newNode()
		rootA := newNode()

		for v := l; v <= r; v++ {
			insert(rootS, v)
		}

		lenA := r - l + 1
		for j := 0; j < lenA; j++ {
			scanner.Scan()
			val, _ := strconv.Atoi(scanner.Text())
			insert(rootA, val)
		}

		pairs := []Pair{{rootS, rootA}}
		ans, _ := solve(MAX_BITS-1, pairs)

		out.WriteString(strconv.Itoa(ans))
		out.WriteByte('\n')
	}
}