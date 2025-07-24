package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 998244353
const maxA = 200000
const K = 18

var spf [maxA + 1]int

func sieve() {
	primes := make([]int, 0)
	for i := 2; i <= maxA; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p > spf[i] || i*p > maxA {
				break
			}
			spf[i*p] = p
		}
	}
}

func factorize(x int) []int {
	res := make([]int, 0, 8)
	for x > 1 {
		p := spf[x]
		res = append(res, p)
		for x%p == 0 {
			x /= p
		}
	}
	return res
}

func main() {
	sieve()

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	primeNodes := make([][]int, maxA+1)
	for i := 1; i <= n; i++ {
		primes := factorize(a[i])
		for _, p := range primes {
			primeNodes[p] = append(primeNodes[p], i)
		}
	}

	up := make([][K + 1]int, n+1)
	depth := make([]int, n+1)
	tin := make([]int, n+1)
	tout := make([]int, n+1)
	timer := 0
	type item struct{ v, p, stage int }
	stack := []item{{1, 1, 0}}
	for len(stack) > 0 {
		it := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		v, p, stage := it.v, it.p, it.stage
		if stage == 0 {
			tin[v] = timer
			timer++
			up[v][0] = p
			for i := 1; i <= K; i++ {
				up[v][i] = up[up[v][i-1]][i-1]
			}
			stack = append(stack, item{v, p, 1})
			for _, w := range adj[v] {
				if w == p {
					continue
				}
				depth[w] = depth[v] + 1
				stack = append(stack, item{w, v, 0})
			}
		} else {
			tout[v] = timer
			timer++
		}
	}

	isAnc := func(u, v int) bool {
		return tin[u] <= tin[v] && tout[v] <= tout[u]
	}
	lca := func(u, v int) int {
		if isAnc(u, v) {
			return u
		}
		if isAnc(v, u) {
			return v
		}
		for i := K; i >= 0; i-- {
			if !isAnc(up[u][i], v) {
				u = up[u][i]
			}
		}
		return up[u][0]
	}

	inv2 := (MOD + 1) / 2
	var ans int64

	for p := 2; p <= maxA; p++ {
		nodes := primeNodes[p]
		m := len(nodes)
		if m < 3 {
			continue
		}
		sort.Slice(nodes, func(i, j int) bool { return tin[nodes[i]] < tin[nodes[j]] })
		all := make([]int, len(nodes))
		copy(all, nodes)
		for i := 0; i < len(nodes)-1; i++ {
			l := lca(nodes[i], nodes[i+1])
			all = append(all, l)
		}
		sort.Slice(all, func(i, j int) bool { return tin[all[i]] < tin[all[j]] })
		uniq := []int{all[0]}
		for i := 1; i < len(all); i++ {
			if all[i] != all[i-1] {
				uniq = append(uniq, all[i])
			}
		}

		vtAdj := make(map[int][]int, len(uniq))
		st := []int{uniq[0]}
		base := make(map[int]bool, m)
		for _, id := range nodes {
			base[id] = true
		}
		for i := 1; i < len(uniq); i++ {
			v := uniq[i]
			l := lca(v, st[len(st)-1])
			for len(st) >= 2 && depth[st[len(st)-2]] >= depth[l] {
				u := st[len(st)-1]
				vtAdj[st[len(st)-2]] = append(vtAdj[st[len(st)-2]], u)
				st = st[:len(st)-1]
			}
			if st[len(st)-1] != l {
				vtAdj[l] = append(vtAdj[l], st[len(st)-1])
				st[len(st)-1] = l
				if len(st) == 1 || st[len(st)-2] != l {
					st = append(st, l)
				}
			}
			st = append(st, v)
		}
		for len(st) > 1 {
			u := st[len(st)-1]
			st = st[:len(st)-1]
			vtAdj[st[len(st)-1]] = append(vtAdj[st[len(st)-1]], u)
		}

		var pairSum int64
		var dfs func(int) int
		dfs = func(u int) int {
			cnt := 0
			if base[u] {
				cnt = 1
			}
			for _, v := range vtAdj[u] {
				c := dfs(v)
				contribution := (int64(c) * int64(m-c)) % MOD
				dist := depth[v] - depth[u]
				contribution = (contribution * int64(dist)) % MOD
				pairSum += contribution
				pairSum %= MOD
				cnt += c
			}
			return cnt
		}
		dfs(uniq[0])
		pairSum %= MOD
		term := pairSum * int64(m-2) % MOD
		term = term * int64(inv2) % MOD
		ans += term
		ans %= MOD
	}

	fmt.Fprintln(out, ans%MOD)
}
