package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const mod int = 1000000007

var sTree [][]int
var tTree [][]int
var fact [13]int

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
	children := []int{}
	for _, v := range tTree[t] {
		if v != pt {
			children = append(children, v)
		}
	}
	neighbors := []int{}
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
		copy(ndp, dp)
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

func expected(nS int, edgesS [][2]int, nT int, edgesT [][2]int) int {
	sTree = make([][]int, nS)
	for _, e := range edgesS {
		u, v := e[0], e[1]
		sTree[u] = append(sTree[u], v)
		sTree[v] = append(sTree[v], u)
	}
	tTree = make([][]int, nT)
	for _, e := range edgesT {
		x, y := e[0], e[1]
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
	return ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// generate random tree with n nodes
func genTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))

	for t := 0; t < 100; t++ {
		nS := rng.Intn(8) + 1
		edgesS := genTree(rng, nS)
		nT := rng.Intn(4) + 1
		edgesT := genTree(rng, nT)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", nS))
		for _, e := range edgesS {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
		}
		sb.WriteString(fmt.Sprintf("%d\n", nT))
		for _, e := range edgesT {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
		}
		input := sb.String()
		exp := fmt.Sprintf("%d", expected(nS, edgesS, nT, edgesT))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\noutput:\n%s\n", t+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
