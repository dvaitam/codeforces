package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type pair struct{ u, v int }

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func expectedAnswerH(t, T int64, li, ri []int64, edges []pair) string {
	n := len(li)
	g := make([][]int, n)
	for _, e := range edges {
		g[e.u] = append(g[e.u], e.v)
		g[e.v] = append(g[e.v], e.u)
	}
	color := make([]int, n)
	comp := make([]int, n)
	var compCol0, compCol1 [][]int
	for i := 0; i < n; i++ {
		if comp[i] != 0 {
			continue
		}
		q := []int{i}
		comp[i] = len(compCol0) + 1
		color[i] = 0
		col0 := []int{i}
		col1 := []int{}
		for k := 0; k < len(q); k++ {
			u := q[k]
			for _, v := range g[u] {
				if comp[v] == 0 {
					comp[v] = comp[u]
					color[v] = color[u] ^ 1
					if color[v] == 0 {
						col0 = append(col0, v)
					} else {
						col1 = append(col1, v)
					}
					q = append(q, v)
				} else {
					if comp[v] == comp[u] && color[v] == color[u] {
						return "IMPOSSIBLE"
					}
				}
			}
		}
		compCol0 = append(compCol0, col0)
		compCol1 = append(compCol1, col1)
	}
	C := len(compCol0)
	type CI struct {
		LA, RA, LB, RB int64
		idx            int
	}
	cis := make([]CI, C)
	for i := 0; i < C; i++ {
		LA, RA := int64(0), int64(1e18)
		for _, u := range compCol0[i] {
			if li[u] > LA {
				LA = li[u]
			}
			if ri[u] < RA {
				RA = ri[u]
			}
		}
		LB, RB := int64(0), int64(1e18)
		for _, u := range compCol1[i] {
			if li[u] > LB {
				LB = li[u]
			}
			if ri[u] < RB {
				RB = ri[u]
			}
		}
		if LA > RA || LB > RB {
			return "IMPOSSIBLE"
		}
		cis[i] = CI{LA, RA, LB, RB, i}
	}
	sort.Slice(cis, func(i, j int) bool { return cis[i].LA-cis[i].LB < cis[j].LA-cis[j].LB })
	const INF = int64(1e18)
	prefLA := make([]int64, C+1)
	prefLB := make([]int64, C+1)
	prefRA := make([]int64, C+1)
	prefRB := make([]int64, C+1)
	prefRA[0], prefRB[0] = INF, INF
	for i := 1; i <= C; i++ {
		c := cis[i-1]
		if prefLA[i-1] > c.LA {
			prefLA[i] = prefLA[i-1]
		} else {
			prefLA[i] = c.LA
		}
		if prefLB[i-1] > c.LB {
			prefLB[i] = prefLB[i-1]
		} else {
			prefLB[i] = c.LB
		}
		if prefRA[i-1] < c.RA {
			prefRA[i] = prefRA[i-1]
		} else {
			prefRA[i] = c.RA
		}
		if prefRB[i-1] < c.RB {
			prefRB[i] = prefRB[i-1]
		} else {
			prefRB[i] = c.RB
		}
	}
	sufLA := make([]int64, C+2)
	sufLB := make([]int64, C+2)
	sufRA := make([]int64, C+2)
	sufRB := make([]int64, C+2)
	sufRA[C+1], sufRB[C+1] = INF, INF
	for i := C; i >= 1; i-- {
		c := cis[i-1]
		if sufLA[i+1] > c.LA {
			sufLA[i] = sufLA[i+1]
		} else {
			sufLA[i] = c.LA
		}
		if sufLB[i+1] > c.LB {
			sufLB[i] = sufLB[i+1]
		} else {
			sufLB[i] = c.LB
		}
		if sufRA[i+1] < c.RA {
			sufRA[i] = sufRA[i+1]
		} else {
			sufRA[i] = c.RA
		}
		if sufRB[i+1] < c.RB {
			sufRB[i] = sufRB[i+1]
		} else {
			sufRB[i] = c.RB
		}
	}
	var bestK int = -1
	var outL1, outL2 int64
	for k := 0; k <= C; k++ {
		L1 := max64(prefLA[k], sufLB[k+1])
		R1 := min64(prefRA[k], sufRB[k+1])
		L2 := max64(prefLB[k], sufLA[k+1])
		R2 := min64(prefRB[k], sufRA[k+1])
		if L1 > R1 || L2 > R2 {
			continue
		}
		if L1+L2 > T || R1+R2 < t {
			continue
		}
		bestK = k
		outL1 = L1
		outL2 = L2
		break
	}
	if bestK < 0 {
		return "IMPOSSIBLE"
	}
	n1 := outL1
	n2 := outL2
	if n1+n2 < t {
		delta := t - (n1 + n2)
		add := min64(delta, min64(prefRA[bestK], sufRB[bestK+1])-n1)
		n1 += add
		delta -= add
		add = min64(delta, min64(prefRB[bestK], sufRA[bestK+1])-n2)
		n2 += add
		delta -= add
	}
	ans := make([]byte, n)
	for i := 0; i < n; i++ {
		ans[i] = '1'
	}
	for idx := 0; idx < bestK; idx++ {
		ci := cis[idx]
		for _, u := range compCol0[ci.idx] {
			ans[u] = '1'
		}
		for _, u := range compCol1[ci.idx] {
			ans[u] = '2'
		}
	}
	for idx := bestK; idx < C; idx++ {
		ci := cis[idx]
		for _, u := range compCol0[ci.idx] {
			ans[u] = '2'
		}
		for _, u := range compCol1[ci.idx] {
			ans[u] = '1'
		}
	}
	return fmt.Sprintf("POSSIBLE\n%d %d\n%s", n1, n2, string(ans))
}

func generateCaseH(rng *rand.Rand) (int64, int64, []int64, []int64, []pair) {
	t := int64(rng.Intn(5) + 1)
	T := t + int64(rng.Intn(5))
	n := rng.Intn(5) + 1
	m := rng.Intn(n)
	li := make([]int64, n)
	ri := make([]int64, n)
	for i := 0; i < n; i++ {
		l := int64(rng.Intn(5))
		r := l + int64(rng.Intn(5))
		li[i] = l
		ri[i] = r
	}
	edges := make([]pair, 0, m)
	seen := map[pair]bool{}
	for len(edges) < m {
		u := rng.Intn(n)
		v := rng.Intn(n)
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		p := pair{u, v}
		if seen[p] {
			continue
		}
		seen[p] = true
		edges = append(edges, p)
	}
	return t, T, li, ri, edges
}

func runCaseH(bin string, t, T int64, li, ri []int64, edges []pair) error {
	var sb strings.Builder
	n := len(li)
	sb.WriteString(fmt.Sprintf("%d %d\n", t, T))
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", li[i], ri[i]))
	}
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u+1, e.v+1))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := strings.TrimSpace(expectedAnswerH(t, T, li, ri, edges))
	if got != expected {
		return fmt.Errorf("expected\n%s\n\ngot\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t, T, li, ri, edges := generateCaseH(rng)
		if err := runCaseH(bin, t, T, li, ri, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
