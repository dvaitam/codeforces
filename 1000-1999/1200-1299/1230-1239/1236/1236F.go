package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

type Edge struct {
	to int
	id int
}

type Cycle struct {
	length int
	ext    int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	adj := make([][]Edge, n+1)
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		edges[i] = [2]int{u, v}
		adj[u] = append(adj[u], Edge{v, i})
		adj[v] = append(adj[v], Edge{u, i})
	}

	deg := make([]int, n+1)
	for i := 1; i <= n; i++ {
		deg[i] = len(adj[i])
	}

	maxPow := 2*n + 5
	inv2 := (MOD + 1) / 2
	powInv2 := make([]int64, maxPow)
	powInv2[0] = 1
	for i := 1; i < maxPow; i++ {
		powInv2[i] = powInv2[i-1] * int64(inv2) % MOD
	}

	visited := make([]bool, n+1)
	visitedEdge := make([]bool, m)
	parent := make([]int, n+1)
	depth := make([]int, n+1)
	type item struct{ v, p, idx int }
	stack := []item{{1, 0, 0}}
	visited[1] = true
	parent[1] = 0
	cycles := make([]Cycle, 0)
	cyclesAt := make([][]int, n+1)

	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		v := top.v
		if top.idx < len(adj[v]) {
			e := adj[v][top.idx]
			top.idx++
			if visitedEdge[e.id] {
				continue
			}
			visitedEdge[e.id] = true
			u := e.to
			if u == top.p {
				continue
			}
			if !visited[u] {
				visited[u] = true
				parent[u] = v
				depth[u] = depth[v] + 1
				stack = append(stack, item{u, v, 0})
			} else if depth[u] < depth[v] {
				// found cycle
				nodes := []int{u}
				ext := deg[u] - 2
				x := v
				for x != u {
					nodes = append(nodes, x)
					ext += deg[x] - 2
					x = parent[x]
				}
				cycleIdx := len(cycles)
				cycles = append(cycles, Cycle{length: len(nodes), ext: ext})
				for _, node := range nodes {
					cyclesAt[node] = append(cyclesAt[node], cycleIdx)
				}
			}
		} else {
			stack = stack[:len(stack)-1]
		}
	}

	c := len(cycles)
	S1 := int64(0)
	S2 := int64(0)
	selfDiff := int64(0)
	for i := 0; i < c; i++ {
		k := cycles[i].length
		pw := powInv2[k]
		S1 = (S1 + pw) % MOD
		S2 = (S2 + powInv2[2*k]) % MOD
		diff := pw - powInv2[2*k]
		if diff < 0 {
			diff += MOD
		}
		selfDiff = (selfDiff + diff) % MOD
	}

	sumShare := int64(0)
	for v := 1; v <= n; v++ {
		lst := cyclesAt[v]
		if len(lst) > 1 {
			var sum, sumSq int64
			for _, idx := range lst {
				k := cycles[idx].length
				sum = (sum + powInv2[k]) % MOD
				sumSq = (sumSq + powInv2[2*k]) % MOD
			}
			contrib := sum*sum%MOD - sumSq
			contrib %= MOD
			if contrib < 0 {
				contrib += MOD
			}
			sumShare = (sumShare + contrib) % MOD
		}
	}

	S1sq := S1 * S1 % MOD
	E_c2 := (S1sq + selfDiff + sumShare) % MOD

	// E[n'^2]
	temp := int64(n)
	E_n2 := (temp*temp%MOD + temp) % MOD
	E_n2 = E_n2 * powInv2[2] % MOD

	// sum of deg choose 2
	sumDegCh2 := int64(0)
	for v := 1; v <= n; v++ {
		d := deg[v]
		sumDegCh2 = (sumDegCh2 + int64(d*(d-1))) % MOD
	}

	// E[m'^2]
	tmpM := int64(m)
	part := (4 * tmpM) % MOD
	part = (part + sumDegCh2) % MOD
	part = (part + tmpM*(tmpM-1)%MOD) % MOD
	E_m2 := part * powInv2[4] % MOD

	// E[n'm']
	E_nm := tmpM * int64((n+2)%MOD) % MOD
	E_nm = E_nm * powInv2[3] % MOD

	// E[n'c']
	E_nc := int64(0)
	for i := 0; i < c; i++ {
		k := cycles[i].length
		val := int64(n+k) % MOD * powInv2[k+1] % MOD
		E_nc = (E_nc + val) % MOD
	}

	// E[m'c']
	E_mc := int64(0)
	for i := 0; i < c; i++ {
		k := cycles[i].length
		ext := cycles[i].ext
		val := (tmpM + int64(3*k) + int64(ext)) % MOD
		val = val * powInv2[k+2] % MOD
		E_mc = (E_mc + val) % MOD
	}

	// E[X^2]
	E_X2 := (E_n2 + E_m2) % MOD
	E_X2 = (E_X2 + E_c2) % MOD
	sub := 2 * E_nm % MOD
	if E_X2 >= sub {
		E_X2 = (E_X2 - sub) % MOD
	} else {
		E_X2 = (E_X2 - sub + MOD) % MOD
	}
	E_X2 = (E_X2 + 2*E_nc%MOD) % MOD
	sub = 2 * E_mc % MOD
	if E_X2 >= sub {
		E_X2 = (E_X2 - sub) % MOD
	} else {
		E_X2 = (E_X2 - sub + MOD) % MOD
	}

	// E[X]
	E_X := temp * powInv2[1] % MOD
	sub = tmpM * powInv2[2] % MOD
	if E_X >= sub {
		E_X = (E_X - sub) % MOD
	} else {
		E_X = (E_X - sub + MOD) % MOD
	}
	E_X = (E_X + S1) % MOD

	E_X_sq := E_X * E_X % MOD
	varX := E_X2 - E_X_sq
	varX %= MOD
	if varX < 0 {
		varX += MOD
	}
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, varX)
	out.Flush()
}
