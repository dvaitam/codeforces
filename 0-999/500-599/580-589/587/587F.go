package main

import (
	"bufio"
	"fmt"
	"os"
)

const B = 320

type Node struct {
	next [26]int
	fail int
}

type BIT struct {
	n    int
	tree []int
}

func newBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+2)}
}

func (b *BIT) add(idx, val int) {
	for idx <= b.n {
		b.tree[idx] += val
		idx += idx & -idx
	}
}

func (b *BIT) rangeAdd(l, r, val int) {
	if l > r {
		return
	}
	b.add(l, val)
	b.add(r+1, -val)
}

func (b *BIT) sum(idx int) int {
	res := 0
	for idx > 0 {
		res += b.tree[idx]
		idx -= idx & -idx
	}
	return res
}

var nodes []Node
var endNode []int
var children [][]int
var tin, tout []int
var order []int
var timer int
var walk [][]int

func newNode() int {
	n := len(nodes)
	var nd Node
	for i := 0; i < 26; i++ {
		nd.next[i] = -1
	}
	nodes = append(nodes, nd)
	return n
}

func insert(s string) int {
	v := 0
	for i := 0; i < len(s); i++ {
		c := int(s[i] - 'a')
		if nodes[v].next[c] == -1 {
			nodes[v].next[c] = newNode()
		}
		v = nodes[v].next[c]
	}
	return v
}

func buildAC() {
	queue := make([]int, 0)
	for c := 0; c < 26; c++ {
		u := nodes[0].next[c]
		if u != -1 {
			nodes[u].fail = 0
			queue = append(queue, u)
		} else {
			nodes[0].next[c] = 0
		}
	}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		f := nodes[v].fail
		for c := 0; c < 26; c++ {
			u := nodes[v].next[c]
			if u != -1 {
				nodes[u].fail = nodes[f].next[c]
				queue = append(queue, u)
			} else {
				nodes[v].next[c] = nodes[f].next[c]
			}
		}
	}
	children = make([][]int, len(nodes))
	for i := 1; i < len(nodes); i++ {
		p := nodes[i].fail
		children[p] = append(children[p], i)
	}
}

func dfs(v int) {
	timer++
	tin[v] = timer
	order = append(order, v)
	for _, u := range children[v] {
		dfs(u)
	}
	tout[v] = timer
}

func querySmall(k int, bit *BIT) int {
	res := 0
	node := 0
	for _, v := range walk[k] {
		node = v
		res += bit.sum(tin[node])
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}

	nodes = make([]Node, 0)
	newNode() // root

	str := make([]string, n+1)
	endNode = make([]int, n+1)
	length := make([]int, n+1)

	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &str[i])
		length[i] = len(str[i])
		endNode[i] = insert(str[i])
	}

	buildAC()

	tin = make([]int, len(nodes))
	tout = make([]int, len(nodes))
	order = make([]int, 0, len(nodes))
	timer = 0
	dfs(0)

	walk = make([][]int, n+1)
	for i := 1; i <= n; i++ {
		path := make([]int, len(str[i]))
		v := 0
		for j := 0; j < len(str[i]); j++ {
			v = nodes[v].next[str[i][j]-'a']
			path[j] = v
		}
		walk[i] = path
	}

	heavyId := make([]int, n+1)
	heavyList := make([]int, 0)
	for i := 1; i <= n; i++ {
		if length[i] > B {
			heavyId[i] = len(heavyList)
			heavyList = append(heavyList, i)
		} else {
			heavyId[i] = -1
		}
	}

	heavyPref := make([][]int, len(heavyList))
	for idx, id := range heavyList {
		freq := make([]int, len(nodes))
		v := 0
		for j := 0; j < len(str[id]); j++ {
			v = nodes[v].next[str[id][j]-'a']
			freq[v]++
		}
		for i := len(order) - 1; i > 0; i-- {
			u := order[i]
			f := nodes[u].fail
			freq[f] += freq[u]
		}
		pref := make([]int, n+1)
		for i := 1; i <= n; i++ {
			pref[i] = pref[i-1] + freq[endNode[i]]
		}
		heavyPref[idx] = pref
	}

	type Event struct {
		idx  int
		k    int
		sign int
	}

	events := make([][]Event, n+1)
	queries := make([][3]int, q)
	ans := make([]int, q)

	for i := 0; i < q; i++ {
		var l, r, k int
		fmt.Fscan(in, &l, &r, &k)
		queries[i] = [3]int{l, r, k}
		if heavyId[k] != -1 {
			hp := heavyPref[heavyId[k]]
			ans[i] = hp[r] - hp[l-1]
		} else {
			events[r] = append(events[r], Event{i, k, 1})
			events[l-1] = append(events[l-1], Event{i, k, -1})
		}
	}

	bit := newBIT(timer + 2)

	for _, ev := range events[0] {
		ans[ev.idx] += ev.sign * querySmall(ev.k, bit)
	}
	for i := 1; i <= n; i++ {
		bit.rangeAdd(tin[endNode[i]], tout[endNode[i]], 1)
		for _, ev := range events[i] {
			ans[ev.idx] += ev.sign * querySmall(ev.k, bit)
		}
	}

	for i := 0; i < q; i++ {
		fmt.Fprintln(out, ans[i])
	}
}
