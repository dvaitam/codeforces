package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

const mod int = 1000000007

var sTree [][]int
var tTree [][]int
var fact [13]int

// memoization key
type state struct {
	u int
	p int
	t int
}

var memo map[state]int

func modMul(a, b int) int {
	return int((int64(a) * int64(b)) % int64(mod))
}

func solve(u, p, t, pt int) int {
	key := state{u, p, t}
	if val, ok := memo[key]; ok {
		return val
	}
	children := make([]int, 0)
	for _, v := range tTree[t] {
		if v != pt {
			children = append(children, v)
		}
	}
	neighbors := make([]int, 0)
	for _, v := range sTree[u] {
		if v != p {
			neighbors = append(neighbors, v)
		}
	}
	if len(neighbors) < len(children) {
		memo[key] = 0
		return 0
	}
	k := len(children)
	m := len(neighbors)
	val := make([][]int, m)
	for j := 0; j < m; j++ {
		val[j] = make([]int, k)
		for i := 0; i < k; i++ {
			val[j][i] = solve(neighbors[j], u, children[i], t)
		}
	}
	dp := make([]int, 1<<k)
	dp[0] = 1
	for j := 0; j < m; j++ {
		ndp := make([]int, 1<<k)
		copy(ndp, dp) // option to not use neighbor
		for mask := 0; mask < (1 << k); mask++ {
			if dp[mask] == 0 {
				continue
			}
			for i := 0; i < k; i++ {
				if (mask>>i)&1 == 0 {
					if val[j][i] == 0 {
						continue
					}
					nm := mask | (1 << i)
					ndp[nm] = (ndp[nm] + modMul(dp[mask], val[j][i])) % mod
				}
			}
		}
		dp = ndp
	}
	res := dp[(1<<k)-1]
	memo[key] = res
	return res
}

func dfsCanon(u, p int) (string, int) {
	forms := []string{}
	aut := 1
	for _, v := range tTree[u] {
		if v == p {
			continue
		}
		s, a := dfsCanon(v, u)
		forms = append(forms, s)
		aut *= a
	}
	sort.Strings(forms)
	i := 0
	for i < len(forms) {
		j := i
		for j < len(forms) && forms[j] == forms[i] {
			j++
		}
		aut *= fact[j-i]
		i = j
	}
	return "(" + strings.Join(forms, "") + ")", aut
}

func modPow(a, e int) int {
	res := 1
	base := a % mod
	for e > 0 {
		if e&1 == 1 {
			res = modMul(res, base)
		}
		base = modMul(base, base)
		e >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var nS int
	if _, err := fmt.Fscan(reader, &nS); err != nil {
		return
	}
	sTree = make([][]int, nS)
	for i := 0; i < nS-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		sTree[u] = append(sTree[u], v)
		sTree[v] = append(sTree[v], u)
	}
	var nT int
	fmt.Fscan(reader, &nT)
	tTree = make([][]int, nT)
	for i := 0; i < nT-1; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		x--
		y--
		tTree[x] = append(tTree[x], y)
		tTree[y] = append(tTree[y], x)
	}
	for i := 0; i < len(fact); i++ {
		if i == 0 {
			fact[i] = 1
		} else {
			fact[i] = fact[i-1] * i
		}
	}
	forms := make([]string, nT)
	autos := make([]int, nT)
	counts := make(map[string]int)
	for v := 0; v < nT; v++ {
		s, a := dfsCanon(v, -1)
		forms[v] = s
		autos[v] = a
		counts[s]++
	}
	rootForm := forms[0]
	autRoot := autos[0]
	orbit := counts[rootForm]
	autTotal := autRoot * orbit

	memo = make(map[state]int)
	total := 0
	for s := 0; s < nS; s++ {
		val := solve(s, -1, 0, -1)
		total += val
		if total >= mod {
			total -= mod
		}
	}
	invAut := modPow(autTotal%mod, mod-2)
	ans := modMul(total, invAut)
	fmt.Println(ans)
}
