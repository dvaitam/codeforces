package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct{ input string }

type edge struct {
	to int
	w  int64
}

func solveCase(in string) string {
	rdr := strings.NewReader(in)
	var n int
	var x int64
	if _, err := fmt.Fscan(rdr, &n, &x); err != nil {
		return ""
	}
	colors := make([]int, n)
	bcnt := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(rdr, &colors[i])
		if colors[i] == 1 {
			bcnt++
		}
	}
	adj := make([][]edge, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		var w int64
		fmt.Fscan(rdr, &u, &v, &w)
		u--
		v--
		adj[u] = append(adj[u], edge{v, w})
		adj[v] = append(adj[v], edge{u, w})
	}
	if bcnt == 0 {
		return "-1"
	}
	dpRoot := dfs(0, -1, x, colors, adj)
	best := -1
	if m, ok := dpRoot[bcnt]; ok {
		for _, v := range m {
			if v > best {
				best = v
			}
		}
	}
	if best < 0 {
		return "-1"
	}
	return fmt.Sprint(bcnt - best)
}

func dfs(u, p int, x int64, color []int, adj [][]edge) map[int]map[int64]int {
	dp0 := map[int]map[int64]int{0: {0: 0}}
	initInter := 0
	if color[u] == 1 {
		initInter = 1
	}
	dp1 := map[int]map[int64]int{1: {-1: initInter}}
	for _, e := range adj[u] {
		v, w := e.to, e.w
		if v == p {
			continue
		}
		dpv := dfs(v, u, x, color, adj)
		new0 := make(map[int]map[int64]int)
		new1 := make(map[int]map[int64]int)
		for k1, m1 := range dp0 {
			for d1, inter1 := range m1 {
				for k2, m2 := range dpv {
					for d2, inter2 := range m2 {
						k := k1 + k2
						var dEff int64
						if d2 < 0 {
							dEff = d1
						} else {
							d2w := d2 + w
							if d2w < d1 {
								dEff = d1
							} else {
								dEff = d2w
							}
						}
						if dEff <= x {
							inter := inter1 + inter2
							mm := new0[k]
							if mm == nil {
								mm = make(map[int64]int)
								new0[k] = mm
							}
							if prev, ok := mm[dEff]; !ok || inter > prev {
								mm[dEff] = inter
							}
						}
					}
				}
			}
		}
		for k1, m1 := range dp1 {
			for d1, inter1 := range m1 {
				for k2, m2 := range dpv {
					for d2, inter2 := range m2 {
						k := k1 + k2
						var dEff int64
						if d2 < 0 {
							dEff = d1
						} else {
							d2w := d2 + w
							if d2w <= x {
								dEff = d1
							} else {
								if d2w < d1 {
									dEff = d1
								} else {
									dEff = d2w
								}
							}
						}
						if dEff <= x {
							inter := inter1 + inter2
							mm := new1[k]
							if mm == nil {
								mm = make(map[int64]int)
								new1[k] = mm
							}
							if prev, ok := mm[dEff]; !ok || inter > prev {
								mm[dEff] = inter
							}
						}
					}
				}
			}
		}
		dp0 = prune(new0)
		dp1 = prune(new1)
	}
	dp := make(map[int]map[int64]int)
	for k, m := range dp0 {
		for d, inter := range m {
			mm := dp[k]
			if mm == nil {
				mm = make(map[int64]int)
				dp[k] = mm
			}
			if prev, ok := mm[d]; !ok || inter > prev {
				mm[d] = inter
			}
		}
	}
	for k, m := range dp1 {
		for d, inter := range m {
			mm := dp[k]
			if mm == nil {
				mm = make(map[int64]int)
				dp[k] = mm
			}
			if prev, ok := mm[d]; !ok || inter > prev {
				mm[d] = inter
			}
		}
	}
	return prune(dp)
}

func prune(dp map[int]map[int64]int) map[int]map[int64]int {
	res := make(map[int]map[int64]int)
	type pair struct {
		d     int64
		inter int
	}
	for k, m := range dp {
		ps := make([]pair, 0, len(m))
		for d, inter := range m {
			ps = append(ps, pair{d, int(inter)})
		}
		for i := 1; i < len(ps); i++ {
			for j := i; j > 0 && (ps[j].d < ps[j-1].d || (ps[j].d == ps[j-1].d && ps[j].inter > ps[j-1].inter)); j-- {
				ps[j], ps[j-1] = ps[j-1], ps[j]
			}
		}
		mm := make(map[int64]int)
		maxInter := -1
		for _, pr := range ps {
			if pr.inter > maxInter {
				mm[pr.d] = pr.inter
				maxInter = pr.inter
			}
		}
		res[k] = mm
	}
	return res
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := solveCase(tc.input)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randomTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 1
	x := int64(rng.Intn(10) + 1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, x))
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteString("0 ")
		} else {
			sb.WriteString("1 ")
		}
	}
	sb.WriteByte('\n')
	edges := randomTree(rng, n)
	for _, e := range edges {
		w := rng.Int63n(10) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], w))
	}
	return testCase{input: sb.String()}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 105)
	cases = append(cases, randomCase(rng))
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
