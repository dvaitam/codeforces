package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1e9 + 7

type Edge struct {
	to    int
	color int
}

type BIT struct {
	tree []int64
}

func (b *BIT) init(n int) {
	b.tree = make([]int64, n+2)
}

func (b *BIT) reset() {
	for i := range b.tree {
		b.tree[i] = 0
	}
}

func (b *BIT) add(i int, v int64) {
	i++
	for i < len(b.tree) {
		b.tree[i] += v
		i += i & -i
	}
}

func (b *BIT) sum(i int) int64 {
	if i < 0 {
		return 0
	}
	if i >= len(b.tree)-1 {
		i = len(b.tree) - 1
	}
	i++
	res := int64(0)
	for i > 0 {
		res += b.tree[i]
		i -= i & -i
	}
	return res
}

type NodeInfo struct {
	cost  int
	sum   int64
	color int
}

var (
	n, k    int
	g       [][]Edge
	a       []int64
	removed []bool
	sz      []int
	ans     int64
	cnt     [2]BIT
	ssum    [2]BIT
)

func dfsSize(u, p int) int {
	sz[u] = 1
	for _, e := range g[u] {
		if e.to != p && !removed[e.to] {
			sz[u] += dfsSize(e.to, u)
		}
	}
	return sz[u]
}

func findCentroid(u, p, total int) int {
	for _, e := range g[u] {
		if e.to != p && !removed[e.to] {
			if sz[e.to]*2 > total {
				return findCentroid(e.to, u, total)
			}
		}
	}
	return u
}

func collect(u, p, firstColor, lastColor, cost int, sum int64, arr *[]NodeInfo) {
	if cost > k {
		return
	}
	sum %= MOD
	*arr = append(*arr, NodeInfo{cost, sum, firstColor})
	for _, e := range g[u] {
		if e.to == p || removed[e.to] {
			continue
		}
		nc := cost
		if lastColor != e.color {
			nc++
		}
		collect(e.to, u, firstColor, e.color, nc, sum+a[e.to], arr)
	}
}

func processCentroid(c int) {
	ans = (ans + a[c]) % MOD
	cnt[0].reset()
	cnt[1].reset()
	ssum[0].reset()
	ssum[1].reset()
	for _, e := range g[c] {
		if removed[e.to] {
			continue
		}
		nodes := make([]NodeInfo, 0)
		collect(e.to, c, e.color, e.color, 0, (a[c]+a[e.to])%MOD, &nodes)
		for _, nd := range nodes {
			if nd.cost <= k {
				ans += nd.sum
				if ans >= MOD {
					ans -= MOD
				}
			}
			diff := 0
			if nd.color != 0 {
				diff = 1
			}
			allowed := k - nd.cost - diff
			if allowed >= 0 {
				cnt0 := cnt[0].sum(allowed)
				sum0 := ssum[0].sum(allowed)
				temp := (cnt0 % MOD) * nd.sum % MOD
				temp = (temp + sum0) % MOD
				temp = (temp - (a[c]%MOD)*(cnt0%MOD)%MOD) % MOD
				if temp < 0 {
					temp += MOD
				}
				ans += temp
				ans %= MOD
			}
			diff = 0
			if nd.color != 1 {
				diff = 1
			}
			allowed = k - nd.cost - diff
			if allowed >= 0 {
				cnt1 := cnt[1].sum(allowed)
				sum1 := ssum[1].sum(allowed)
				temp := (cnt1 % MOD) * nd.sum % MOD
				temp = (temp + sum1) % MOD
				temp = (temp - (a[c]%MOD)*(cnt1%MOD)%MOD) % MOD
				if temp < 0 {
					temp += MOD
				}
				ans += temp
				ans %= MOD
			}
		}
		for _, nd := range nodes {
			if nd.cost <= k {
				if nd.color == 0 {
					cnt[0].add(nd.cost, 1)
					ssum[0].add(nd.cost, nd.sum)
				} else {
					cnt[1].add(nd.cost, 1)
					ssum[1].add(nd.cost, nd.sum)
				}
			}
		}
	}
}

func decompose(u int) {
	total := dfsSize(u, -1)
	c := findCentroid(u, -1, total)
	processCentroid(c)
	removed[c] = true
	for _, e := range g[c] {
		if !removed[e.to] {
			decompose(e.to)
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &k)
	a = make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		a[i] %= MOD
	}
	g = make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		var u, v, t int
		fmt.Fscan(in, &u, &v, &t)
		u--
		v--
		g[u] = append(g[u], Edge{v, t})
		g[v] = append(g[v], Edge{u, t})
	}
	removed = make([]bool, n)
	sz = make([]int, n)
	for i := 0; i < 2; i++ {
		cnt[i].init(k + 2)
		ssum[i].init(k + 2)
	}
	decompose(0)
	ans %= MOD
	if ans < 0 {
		ans += MOD
	}
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans%MOD)
	out.Flush()
}
