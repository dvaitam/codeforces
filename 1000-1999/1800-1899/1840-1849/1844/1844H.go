package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

// DSU structure tracking number of vertices and edges in each component
// to detect when a cycle is formed.
type DSU struct {
	parent []int
	size   []int
	edges  []int
}

func NewDSU(n int) *DSU {
	d := &DSU{
		parent: make([]int, n+1),
		size:   make([]int, n+1),
		edges:  make([]int, n+1),
	}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) addEdge(u, v int) (root, edges, size int) {
	ru := d.find(u)
	rv := d.find(v)
	if ru != rv {
		d.parent[ru] = rv
		d.size[rv] += d.size[ru]
		d.edges[rv] += d.edges[ru] + 1
		root = rv
	} else {
		d.edges[ru]++
		root = ru
	}
	root = d.find(root)
	edges = d.edges[root]
	size = d.size[root]
	return
}

// precompute factorials and helper arrays for permutations with cycle lengths multiples of 3
func precompute(n int) (fact, invFact, powInv3, prefProd []int64) {
	fact = make([]int64, n+1)
	invFact = make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[n] = powMod(fact[n], MOD-2)
	for i := n; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
	mMax := n / 3
	powInv3 = make([]int64, mMax+1)
	prefProd = make([]int64, mMax+1)
	powInv3[0] = 1
	prefProd[0] = 1
	inv3 := powMod(3, MOD-2)
	for i := 1; i <= mMax; i++ {
		powInv3[i] = powInv3[i-1] * inv3 % MOD
		prefProd[i] = prefProd[i-1] * int64(3*(i-1)+1) % MOD
	}
	return
}

func powMod(a int64, e int64) int64 {
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

func waysRemaining(r int, fact, invFact, powInv3, prefProd []int64) int64 {
	if r%3 != 0 {
		return 0
	}
	m := r / 3
	ans := fact[r]
	ans = ans * prefProd[m] % MOD
	ans = ans * powInv3[m] % MOD
	ans = ans * invFact[m] % MOD
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	fact, invFact, powInv3, prefProd := precompute(n)

	dsu := NewDSU(n)
	invalid := false
	for t := 1; t <= n; t++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		if !invalid {
			_, edges, size := dsu.addEdge(x, y)
			if edges == size && size%3 != 0 {
				invalid = true
			}
		}
		r := n - t
		if invalid {
			fmt.Fprintln(out, 0)
		} else {
			fmt.Fprintln(out, waysRemaining(r, fact, invFact, powInv3, prefProd))
		}
	}
}
