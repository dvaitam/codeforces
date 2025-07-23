package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT { return &BIT{n: n, tree: make([]int, n+2)} }
func (b *BIT) Add(i, val int) {
	for i <= b.n {
		b.tree[i] += val
		i += i & -i
	}
}
func (b *BIT) Sum(i int) int {
	s := 0
	for i > 0 {
		s += b.tree[i]
		i -= i & -i
	}
	return s
}
func (b *BIT) Range(l, r int) int {
	if l > r {
		return 0
	}
	return b.Sum(r) - b.Sum(l-1)
}

type NodeExp struct{ node, exp int }
type QueryPrime struct{ idx, u, v, exp int }

const MOD int64 = 1000000007

var (
	g      [][]int
	parent []int
	depth  []int
	heavy  []int
	sizeT  []int
	top    []int
	pos    []int
	curPos int
)

func dfs(u, p int) {
	parent[u] = p
	sizeT[u] = 1
	heavy[u] = -1
	for _, v := range g[u] {
		if v == p {
			continue
		}
		depth[v] = depth[u] + 1
		dfs(v, u)
		sizeT[u] += sizeT[v]
		if heavy[u] == -1 || sizeT[v] > sizeT[heavy[u]] {
			heavy[u] = v
		}
	}
}

func decompose(u, h int) {
	top[u] = h
	pos[u] = curPos
	curPos++
	if heavy[u] != -1 {
		decompose(heavy[u], h)
		for _, v := range g[u] {
			if v != parent[u] && v != heavy[u] {
				decompose(v, v)
			}
		}
	}
}

func queryPath(u, v int, bit *BIT) int {
	res := 0
	for top[u] != top[v] {
		if depth[top[u]] < depth[top[v]] {
			u, v = v, u
		}
		res += bit.Range(pos[top[u]]+1, pos[u]+1)
		u = parent[top[u]]
	}
	if depth[u] > depth[v] {
		u, v = v, u
	}
	res += bit.Range(pos[u]+1, pos[v]+1)
	return res
}

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

func factorize(x int, spf []int) map[int]int {
	m := make(map[int]int)
	for x > 1 {
		p := spf[x]
		c := 0
		for x%p == 0 {
			x /= p
			c++
		}
		m[p] += c
	}
	return m
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	g = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	var q int
	fmt.Fscan(in, &q)
	qs := make([]struct{ u, v, x int }, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &qs[i].u, &qs[i].v, &qs[i].x)
	}

	// sieve for spf
	maxV := 10000000
	spf := make([]int, maxV+1)
	primes := make([]int, 0)
	for i := 2; i <= maxV; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p > spf[i] || i*p > maxV {
				break
			}
			spf[i*p] = p
		}
	}

	// HLD preprocess
	parent = make([]int, n+1)
	depth = make([]int, n+1)
	heavy = make([]int, n+1)
	sizeT = make([]int, n+1)
	top = make([]int, n+1)
	pos = make([]int, n+1)
	dfs(1, 0)
	curPos = 0
	decompose(1, 1)

	// factor nodes and queries per prime
	primeNodes := make(map[int][]NodeExp)
	for v := 1; v <= n; v++ {
		m := factorize(a[v], spf)
		for p, e := range m {
			primeNodes[p] = append(primeNodes[p], NodeExp{v, e})
		}
	}
	primeQueries := make(map[int][]QueryPrime)
	for idx := 0; idx < q; idx++ {
		m := factorize(qs[idx].x, spf)
		for p, e := range m {
			primeQueries[p] = append(primeQueries[p], QueryPrime{idx, qs[idx].u, qs[idx].v, e})
		}
	}

	answers := make([]int64, q)
	for i := 0; i < q; i++ {
		answers[i] = 1
	}
	bit := NewBIT(n)

	// iterate over primes appearing in queries
	for p, qsList := range primeQueries {
		nodes := primeNodes[p]
		if len(nodes) == 0 && len(qsList) == 0 {
			continue
		}
		// find max exponent
		maxE := 0
		for _, ne := range nodes {
			if ne.exp > maxE {
				maxE = ne.exp
			}
		}
		for _, qu := range qsList {
			if qu.exp > maxE {
				maxE = qu.exp
			}
		}
		// build levels
		levels := make([][]QueryPrime, maxE+1)
		for _, qu := range qsList {
			for t := 1; t <= qu.exp; t++ {
				levels[t] = append(levels[t], qu)
			}
		}
		sort.Slice(nodes, func(i, j int) bool { return nodes[i].exp > nodes[j].exp })
		ptr := 0
		for t := maxE; t >= 1; t-- {
			for ptr < len(nodes) && nodes[ptr].exp >= t {
				bit.Add(pos[nodes[ptr].node]+1, 1)
				ptr++
			}
			for _, qu := range levels[t] {
				cnt := queryPath(qu.u, qu.v, bit)
				answers[qu.idx] = answers[qu.idx] * modPow(int64(p), int64(cnt)) % MOD
			}
		}
		for i := 0; i < ptr; i++ {
			bit.Add(pos[nodes[i].node]+1, -1)
		}
	}

	for i := 0; i < q; i++ {
		fmt.Fprintln(out, answers[i]%MOD)
	}
}
