package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 998244353
const MAXN = 5000

var catalan []int64

type Pair struct {
	l, r int
}

type Node struct {
	l, r     int
	children []*Node
}

func add(a, b int64) int64 {
	res := (a + b) % MOD
	if res < 0 {
		res += MOD
	}
	return res
}

func mul(a, b int64) int64 {
	return (a * b) % MOD
}

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = mul(res, a)
		}
		a = mul(a, a)
		e >>= 1
	}
	return res
}

func prepareCatalan() {
	limit := 2 * MAXN
	fact := make([]int64, limit+1)
	ifact := make([]int64, limit+1)
	fact[0] = 1
	for i := 1; i <= limit; i++ {
		fact[i] = mul(fact[i-1], int64(i))
	}
	ifact[limit] = modPow(fact[limit], MOD-2)
	for i := limit; i >= 1; i-- {
		ifact[i-1] = mul(ifact[i], int64(i))
	}
	catalan = make([]int64, MAXN+1)
	for n := 0; n <= MAXN; n++ {
		num := fact[2*n]
		den := mul(ifact[n], ifact[n])
		c := mul(num, den)
		c = mul(c, modPow(int64(n+1), MOD-2))
		catalan[n] = c
	}
}

func countWays(n int, pairs []Pair) int64 {
	if len(pairs) == 0 {
		return catalan[n]
	}
	arr := append([]Pair(nil), pairs...)
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].l == arr[j].l {
			return arr[i].r < arr[j].r
		}
		return arr[i].l < arr[j].l
	})

	root := &Node{l: 0, r: 2*n + 1}
	stack := []*Node{root}

	for _, p := range arr {
		if p.r <= p.l {
			return 0
		}
		for len(stack) > 0 {
			top := stack[len(stack)-1]
			if top.l < p.l && p.r < top.r {
				node := &Node{l: p.l, r: p.r}
				top.children = append(top.children, node)
				stack = append(stack, node)
				break
			}
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			return 0
		}
	}

	var dfs func(*Node) int64
	dfs = func(node *Node) int64 {
		res := int64(1)
		total := node.r - node.l - 1
		used := 0
		for _, ch := range node.children {
			res = mul(res, dfs(ch))
			length := ch.r - ch.l + 1
			used += length
		}
		remaining := total - used
		if remaining < 0 || remaining%2 != 0 {
			return 0
		}
		res = mul(res, catalan[remaining/2])
		return res
	}

	return dfs(root)
}

func main() {
	prepareCatalan()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		pairs := make([]Pair, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &pairs[i].l, &pairs[i].r)
		}
		results := make([]int64, n+1)
		results[0] = catalan[n]
		for i := 1; i <= n; i++ {
			results[i] = countWays(n, pairs[:i])
		}
		for i, v := range results {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
