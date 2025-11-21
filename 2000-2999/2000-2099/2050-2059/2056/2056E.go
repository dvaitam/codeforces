package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 998244353
const MAXN = 400000 + 5

var fact [MAXN]int64
var invFact [MAXN]int64

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func init() {
	fact[0] = 1
	for i := 1; i < MAXN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[MAXN-1] = modPow(fact[MAXN-1], MOD-2)
	for i := MAXN - 2; i >= 0; i-- {
		invFact[i] = invFact[i+1] * int64(i+1) % MOD
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
}

func catalan(k int) int64 {
	if k < 0 {
		return 0
	}
	c := comb(2*k, k)
	inv := modPow(int64(k+1), MOD-2)
	return c * inv % MOD
}

type Node struct {
	l, r     int
	children []int
}

var nodes []Node

func dfs(idx int) int64 {
	node := &nodes[idx]
	if node.l == node.r {
		return 1
	}
	sort.Slice(node.children, func(i, j int) bool {
		return nodes[node.children[i]].l < nodes[node.children[j]].l
	})
	cur := node.l
	blocks := 0
	res := int64(1)
	for _, ch := range node.children {
		child := &nodes[ch]
		if cur < child.l {
			blocks += child.l - cur
			cur = child.l
		}
		res = res * dfs(ch) % MOD
		blocks++
		cur = child.r + 1
	}
	if cur <= node.r {
		blocks += node.r - cur + 1
	}
	res = res * catalan(blocks-1) % MOD
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		nodes = make([]Node, m+1)
		segs := make([]Node, 0, m+1)
		segs = append(segs, Node{l: 1, r: n})
		for i := 0; i < m; i++ {
			var l, r int
			fmt.Fscan(in, &l, &r)
			segs = append(segs, Node{l: l, r: r})
		}
		sort.Slice(segs, func(i, j int) bool {
			if segs[i].l == segs[j].l {
				return segs[i].r > segs[j].r
			}
			return segs[i].l < segs[j].l
		})
		nodes = make([]Node, len(segs))
		for i, seg := range segs {
			nodes[i] = Node{l: seg.l, r: seg.r}
		}
		stack := make([]int, 0)
		for i := range nodes {
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if nodes[top].l <= nodes[i].l && nodes[i].r <= nodes[top].r {
					break
				}
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				nodes[parent].children = append(nodes[parent].children, i)
			}
			stack = append(stack, i)
		}
		ans := dfs(0)
		fmt.Fprintln(out, ans)
	}
}
