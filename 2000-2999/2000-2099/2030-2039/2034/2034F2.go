package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod = 998244353

var fact, invFact []int

func modPow(a, e int) int {
	res := 1
	base := a
	exp := e
	for exp > 0 {
		if exp&1 == 1 {
			res = int(int64(res) * int64(base) % mod)
		}
		base = int(int64(base) * int64(base) % mod)
		exp >>= 1
	}
	return res
}

func modInv(a int) int {
	return modPow(a, mod-2)
}

func comb(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	return int(int64(fact[n]) * int64(invFact[k]) % mod * int64(invFact[n-k]) % mod)
}

type node struct {
	r, b  int
	spec  bool
	index int
	sum   int
}

func prepareFacts(limit int) {
	if len(fact) >= limit+1 {
		return
	}
	fact = make([]int, limit+1)
	invFact = make([]int, limit+1)
	fact[0] = 1
	for i := 1; i <= limit; i++ {
		fact[i] = int(int64(fact[i-1]) * int64(i) % mod)
	}
	invFact[limit] = modInv(fact[limit])
	for i := limit; i > 0; i-- {
		invFact[i-1] = int(int64(invFact[i]) * int64(i) % mod)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)

	// Gather maximum n+m to size factorials once.
	type testCase struct {
		n, m     int
		specials [][2]int
	}
	tests := make([]testCase, T)
	maxSum := 0
	for t := 0; t < T; t++ {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		tests[t].n = n
		tests[t].m = m
		tests[t].specials = make([][2]int, k)
		if n+m > maxSum {
			maxSum = n + m
		}
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &tests[t].specials[i][0], &tests[t].specials[i][1])
		}
	}

	prepareFacts(maxSum)

	for _, tc := range tests {
		n, m := tc.n, tc.m
		k := len(tc.specials)
		nodes := make([]node, 0, k+2)
		nodes = append(nodes, node{r: n, b: m, spec: false, sum: n + m})
		for i := 0; i < k; i++ {
			r := tc.specials[i][0]
			b := tc.specials[i][1]
			nodes = append(nodes, node{r: r, b: b, spec: true, sum: r + b})
		}
		nodes = append(nodes, node{r: 0, b: 0, spec: false, sum: 0})

		// Sort by remaining total descending, tie by red descending to ensure deterministic order.
		sort.Slice(nodes, func(i, j int) bool {
			if nodes[i].sum != nodes[j].sum {
				return nodes[i].sum > nodes[j].sum
			}
			return nodes[i].r > nodes[j].r
		})

		N := len(nodes)
		// Recompute index after sort if needed.
		for i := range nodes {
			nodes[i].index = i
		}

		// Build reachability list.
		reach := make([][]int, N)
		for i := 0; i < N; i++ {
			for j := i + 1; j < N; j++ {
				if nodes[i].r >= nodes[j].r && nodes[i].b >= nodes[j].b {
					reach[i] = append(reach[i], j)
				}
			}
		}

		// direct[i] aligns with reach[i]
		direct := make([][]int, N)
		for i := 0; i < N; i++ {
			adj := reach[i]
			if len(adj) == 0 {
				continue
			}
			direct[i] = make([]int, len(adj))
			for p, j := range adj {
				dr := nodes[i].r - nodes[j].r
				db := nodes[i].b - nodes[j].b
				val := comb(dr+db, dr)
				for q := 0; q < p; q++ {
					kidx := adj[q]
					// Check reachability for safety.
					if nodes[kidx].r < nodes[j].r || nodes[kidx].b < nodes[j].b {
						continue
					}
					dr2 := nodes[kidx].r - nodes[j].r
					db2 := nodes[kidx].b - nodes[j].b
					ways := comb(dr2+db2, dr2)
					val = int((int64(val) - int64(direct[i][q])*int64(ways)) % mod)
					if val < 0 {
						val += mod
					}
				}
				direct[i][p] = val
			}
		}

		paths := make([]int, N)
		vals := make([]int, N)
		paths[0] = 1
		vals[0] = 0

		for i := 0; i < N; i++ {
			if paths[i] == 0 {
				continue
			}
			adj := reach[i]
			dirRow := direct[i]
			for idx, j := range adj {
				d := dirRow[idx]
				if d == 0 {
					continue
				}
				dr := nodes[i].r - nodes[j].r
				db := nodes[i].b - nodes[j].b
				add := (2*dr + db) % mod

				cp := int(int64(paths[i]) * int64(d) % mod)
				pre := (int64(vals[i]) + int64(paths[i])*int64(add)) % mod
				cv := int(pre * int64(d) % mod)
				if nodes[j].spec {
					cv = (cv * 2) % mod
				}

				paths[j] += cp
				if paths[j] >= mod {
					paths[j] -= mod
				}
				vals[j] += cv
				if vals[j] >= mod {
					vals[j] -= mod
				}
			}
		}

		totalPaths := comb(n+m, n)
		invPaths := modInv(totalPaths)
		ans := int(int64(vals[N-1]) * int64(invPaths) % mod)
		fmt.Fprintln(out, ans)
	}
}
