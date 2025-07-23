package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

var (
	n           int
	a           []int
	tree        [][]int
	phi         []int
	phiWeight   []int64
	divisors    [][]int
	nodeDivs    [][]int
	resMultiple []int64
	countArr    []int64
	distArr     []int64
	usedDivs    []int
	removed     []bool
	sz          []int
)

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

func sievePhi(m int) []int {
	phi := make([]int, m+1)
	for i := 0; i <= m; i++ {
		phi[i] = i
	}
	for i := 2; i <= m; i++ {
		if phi[i] == i {
			for j := i; j <= m; j += i {
				phi[j] = phi[j] / i * (i - 1)
			}
		}
	}
	return phi
}

func calcDivs(m int) [][]int {
	d := make([][]int, m+1)
	for i := 1; i <= m; i++ {
		for j := i; j <= m; j += i {
			d[j] = append(d[j], i)
		}
	}
	return d
}

func dfsSize(v, p int) int {
	sz[v] = 1
	for _, to := range tree[v] {
		if to != p && !removed[to] {
			sz[v] += dfsSize(to, v)
		}
	}
	return sz[v]
}

func findCentroid(v, p, tot int) int {
	for _, to := range tree[v] {
		if to != p && !removed[to] && sz[to] > tot/2 {
			return findCentroid(to, v, tot)
		}
	}
	return v
}

func collect(v, p, depth int, nodes *[]int, depths *[]int) {
	*nodes = append(*nodes, v)
	*depths = append(*depths, depth)
	for _, to := range tree[v] {
		if to != p && !removed[to] {
			collect(to, v, depth+1, nodes, depths)
		}
	}
}

func addNode(v, depth int) {
	w := phiWeight[v]
	for _, d := range nodeDivs[v] {
		if countArr[d] == 0 && distArr[d] == 0 {
			usedDivs = append(usedDivs, d)
		}
		countArr[d] = (countArr[d] + w) % MOD
		distArr[d] = (distArr[d] + w*int64(depth)) % MOD
	}
}

func queryNode(v, depth int) {
	w := phiWeight[v]
	for _, d := range nodeDivs[v] {
		val := (distArr[d] + int64(depth)*countArr[d]) % MOD
		resMultiple[d] = (resMultiple[d] + w*val) % MOD
	}
}

func decompose(v int) {
	tot := dfsSize(v, -1)
	c := findCentroid(v, -1, tot)
	removed[c] = true
	usedDivs = usedDivs[:0]
	addNode(c, 0)
	for _, to := range tree[c] {
		if !removed[to] {
			nodes := make([]int, 0)
			depths := make([]int, 0)
			collect(to, c, 1, &nodes, &depths)
			for i, node := range nodes {
				queryNode(node, depths[i])
			}
			for i, node := range nodes {
				addNode(node, depths[i])
			}
		}
	}
	for _, d := range usedDivs {
		countArr[d] = 0
		distArr[d] = 0
	}
	for _, to := range tree[c] {
		if !removed[to] {
			decompose(to)
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n)
	a = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	tree = make([][]int, n)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		tree[x] = append(tree[x], y)
		tree[y] = append(tree[y], x)
	}

	phi = sievePhi(n)
	divisors = calcDivs(n)
	phiWeight = make([]int64, n)
	nodeDivs = make([][]int, n)
	for i := 0; i < n; i++ {
		phiWeight[i] = int64(phi[a[i]])
		nodeDivs[i] = divisors[a[i]]
	}
	resMultiple = make([]int64, n+1)
	countArr = make([]int64, n+1)
	distArr = make([]int64, n+1)
	usedDivs = make([]int, 0)
	removed = make([]bool, n)
	sz = make([]int, n)

	decompose(0)

	resExact := make([]int64, n+1)
	for i := n; i >= 1; i-- {
		val := resMultiple[i] % MOD
		for j := i * 2; j <= n; j += i {
			val = (val - resExact[j]) % MOD
		}
		if val < 0 {
			val += MOD
		}
		resExact[i] = val
	}

	numerator := int64(0)
	for g := 1; g <= n; g++ {
		if resExact[g] == 0 {
			continue
		}
		factor := int64(g) * modPow(int64(phi[g]), MOD-2) % MOD
		numerator = (numerator + resExact[g]*factor) % MOD
	}
	numerator = numerator * 2 % MOD
	denom := int64(n) * int64(n-1) % MOD
	ans := numerator * modPow(denom, MOD-2) % MOD
	fmt.Println(ans)
}
