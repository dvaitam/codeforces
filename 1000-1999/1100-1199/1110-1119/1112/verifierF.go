package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n int
	c []int64
	g [][]int
}

type dsu struct {
	p []int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return &dsu{p: p}
}

func (d *dsu) find(x int) int {
	if d.p[x] == x {
		return x
	}
	d.p[x] = d.find(d.p[x])
	return d.p[x]
}

func (d *dsu) union(a, b int) {
	ra, rb := d.find(a), d.find(b)
	if ra != rb {
		d.p[ra] = rb
	}
}

func genCase(r *rand.Rand) string {
	// Problem constraints: n >= 2
	n := r.Intn(5) + 2
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		// Problem constraints: 0 <= c_i <= 1e9
		sb.WriteString(fmt.Sprint(r.Intn(11)))
	}
	sb.WriteByte('\n')
	for i := 2; i <= n; i++ {
		v := r.Intn(i-1) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", v, i))
	}
	return sb.String()
}

func parseCase(input string) (testCase, error) {
	in := strings.NewReader(input)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return testCase{}, err
	}
	c := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		if _, err := fmt.Fscan(in, &c[i]); err != nil {
			return testCase{}, err
		}
	}
	g := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var a, b int
		if _, err := fmt.Fscan(in, &a, &b); err != nil {
			return testCase{}, err
		}
		g[a] = append(g[a], b)
		g[b] = append(g[b], a)
	}
	return testCase{n: n, c: c, g: g}, nil
}

func buildTinTout(tc testCase) ([]int, []int, int) {
	tin := make([]int, tc.n+1)
	tout := make([]int, tc.n+1)
	tot := 0
	var dfs func(u, p int)
	dfs = func(u, p int) {
		tin[u] = tot
		if u != 1 && len(tc.g[u]) == 1 {
			tot++
		}
		for _, v := range tc.g[u] {
			if v == p {
				continue
			}
			dfs(v, u)
		}
		tout[u] = tot
	}
	dfs(1, 0)
	return tin, tout, tot
}

func isWinning(mask int, tc testCase, tin, tout []int, leafCnt int) bool {
	d := newDSU(leafCnt + 1)
	for v := 1; v <= tc.n; v++ {
		if (mask&(1<<(v-1))) == 0 || tin[v] == tout[v] {
			continue
		}
		d.union(tin[v], tout[v])
	}
	root := d.find(0)
	for x := 1; x <= leafCnt; x++ {
		if d.find(x) != root {
			return false
		}
	}
	return true
}

func solveBruteforce(tc testCase) (int64, []int) {
	tin, tout, leafCnt := buildTinTout(tc)
	const inf int64 = 1<<62 - 1
	best := inf
	totalMasks := 1 << tc.n
	for mask := 1; mask < totalMasks; mask++ {
		if !isWinning(mask, tc, tin, tout, leafCnt) {
			continue
		}
		var sum int64
		for v := 1; v <= tc.n; v++ {
			if (mask & (1 << (v - 1))) != 0 {
				sum += tc.c[v]
			}
		}
		if sum < best {
			best = sum
		}
	}

	good := make([]bool, tc.n+1)
	for mask := 1; mask < totalMasks; mask++ {
		if !isWinning(mask, tc, tin, tout, leafCnt) {
			continue
		}
		var sum int64
		for v := 1; v <= tc.n; v++ {
			if (mask & (1 << (v - 1))) != 0 {
				sum += tc.c[v]
			}
		}
		if sum != best {
			continue
		}
		for v := 1; v <= tc.n; v++ {
			if (mask & (1 << (v - 1))) != 0 {
				good[v] = true
			}
		}
	}

	res := make([]int, 0)
	for v := 1; v <= tc.n; v++ {
		if good[v] {
			res = append(res, v)
		}
	}
	return best, res
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseOutput(out string, n int) (int64, []int, error) {
	fields := strings.Fields(out)
	if len(fields) < 2 {
		return 0, nil, fmt.Errorf("output too short: %q", out)
	}
	s, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid sum %q", fields[0])
	}
	k64, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil || k64 < 0 {
		return 0, nil, fmt.Errorf("invalid k %q", fields[1])
	}
	k := int(k64)
	if len(fields) != 2+k {
		return 0, nil, fmt.Errorf("expected %d vertex indices, got %d tokens", k, len(fields)-2)
	}
	seen := make(map[int]bool, k)
	verts := make([]int, 0, k)
	for i := 0; i < k; i++ {
		v64, err := strconv.ParseInt(fields[2+i], 10, 64)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid vertex index %q", fields[2+i])
		}
		v := int(v64)
		if v < 1 || v > n {
			return 0, nil, fmt.Errorf("vertex index out of range: %d", v)
		}
		if seen[v] {
			return 0, nil, fmt.Errorf("duplicate vertex index: %d", v)
		}
		seen[v] = true
		verts = append(verts, v)
	}
	sort.Ints(verts)
	return s, verts, nil
}

func equalIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func formatAnswer(s int64, verts []int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", s, len(verts)))
	for i, v := range verts {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		tc, err := parseCase(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal parse error on case %d: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		expectS, expectVerts := solveBruteforce(tc)

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		gotS, gotVerts, err := parseOutput(got, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid output: %v\ngot:\n%s\ninput:\n%s", i, err, got, input)
			os.Exit(1)
		}
		if gotS != expectS || !equalIntSlices(gotVerts, expectVerts) {
			fmt.Fprintf(
				os.Stderr,
				"case %d failed: expected %s got %s\ninput:\n%s",
				i,
				formatAnswer(expectS, expectVerts),
				formatAnswer(gotS, gotVerts),
				input,
			)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
