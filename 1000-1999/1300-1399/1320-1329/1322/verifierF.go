package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n     int
	m     int
	edges [][2]int
	lines [][2]int
}

type pair struct {
	l int
	r int
}

const testcasesRaw = `4 2 1 2 2 3 2 4 4 2 1 3
2 3 1 2 1 2 1 2 1 2
4 1 1 2 2 3 1 4 4 3
3 3 1 2 1 3 3 2 2 3 3 1
2 3 1 2 2 1 1 2 1 2
5 3 1 2 1 3 2 4 4 5 1 2 2 3 3 2
3 1 1 2 2 3 3 1
3 3 1 2 1 3 3 2 2 1 2 1
5 3 1 2 2 3 1 4 1 5 5 4 4 5 3 1
2 1 1 2 2 1
2 2 1 2 2 1 1 2
4 2 1 2 2 3 3 4 3 1 3 4
2 3 1 2 2 1 1 2 1 2
5 2 1 2 2 3 2 4 1 5 2 3 3 4
4 1 1 2 1 3 1 4 3 2
2 2 1 2 1 2 1 2
5 3 1 2 1 3 2 4 4 5 3 4 3 2 1 4
3 2 1 2 1 3 1 3 1 2
4 3 1 2 2 3 3 4 4 2 4 1 2 4
5 3 1 2 1 3 1 4 1 5 3 5 4 2 2 5
5 3 1 2 1 3 3 4 1 5 2 1 5 1 2 5
2 1 1 2 1 2
5 3 1 2 1 3 2 4 1 5 3 4 1 5 3 1
2 1 1 2 2 1
3 1 1 2 2 3 2 1
5 3 1 2 1 3 3 4 2 5 1 4 2 4 5 4
3 1 1 2 1 3 2 1
3 2 1 2 1 3 3 2 1 2
4 2 1 2 2 3 3 4 1 3 3 4
4 2 1 2 1 3 1 4 4 2 1 2
4 1 1 2 1 3 3 4 2 3
5 2 1 2 2 3 3 4 4 5 2 1 2 3
3 2 1 2 1 3 3 2 2 3
4 1 1 2 2 3 1 4 2 1
4 1 1 2 1 3 3 4 1 2
5 3 1 2 1 3 2 4 1 5 2 5 5 1 2 1
5 2 1 2 1 3 3 4 4 5 3 1 2 1
3 2 1 2 1 3 3 1 2 1
2 2 1 2 2 1 2 1
5 1 1 2 2 3 1 4 1 5 3 5
3 3 1 2 1 3 1 3 2 1 3 1
4 2 1 2 1 3 3 4 4 3 3 4
4 2 1 2 1 3 1 4 1 2 1 3
4 3 1 2 2 3 1 4 4 1 3 4 4 3
2 3 1 2 2 1 2 1 1 2
5 2 1 2 1 3 2 4 2 5 5 3 2 1
4 1 1 2 2 3 3 4 4 3
5 1 1 2 1 3 1 4 1 5 5 4
2 2 1 2 1 2 2 1
4 2 1 2 2 3 2 4 3 2 4 2
2 1 1 2 2 1
4 1 1 2 1 3 1 4 2 3
4 3 1 2 2 3 3 4 1 3 2 4 2 3
2 1 1 2 2 1
2 1 1 2 2 1
4 3 1 2 2 3 2 4 2 3 3 2 3 2
2 2 1 2 2 1 1 2
3 3 1 2 2 3 1 3 2 3 3 2
4 1 1 2 2 3 1 4 1 4
5 1 1 2 2 3 1 4 1 5 4 2
4 1 1 2 2 3 1 4 3 4
4 2 1 2 2 3 2 4 3 1 1 2
5 2 1 2 2 3 3 4 1 5 5 2 2 4
4 1 1 2 1 3 1 4 4 2
2 2 1 2 2 1 1 2
2 3 1 2 2 1 1 2 2 1
4 3 1 2 1 3 2 4 3 4 3 1 4 2
4 2 1 2 1 3 1 4 1 3 4 2
4 1 1 2 1 3 2 4 4 3
3 2 1 2 2 3 2 3 3 2
5 2 1 2 1 3 3 4 4 5 3 2 1 5
3 1 1 2 2 3 2 1
4 1 1 2 1 3 1 4 3 1
4 3 1 2 2 3 1 4 2 3 2 3 2 1
2 1 1 2 1 2
2 2 1 2 2 1 1 2
5 1 1 2 2 3 3 4 4 5 2 4
2 1 1 2 2 1
5 2 1 2 2 3 2 4 2 5 4 5 4 5
3 3 1 2 1 3 1 3 3 2 3 1
4 2 1 2 2 3 3 4 3 1 1 3
3 2 1 2 1 3 1 3 1 2
5 2 1 2 1 3 1 4 2 5 2 1 4 2
4 1 1 2 1 3 2 4 1 2
5 1 1 2 2 3 1 4 2 5 1 3
5 3 1 2 2 3 1 4 4 5 4 5 2 4 1 5
2 3 1 2 2 1 1 2 2 1
2 2 1 2 1 2 2 1
4 2 1 2 2 3 1 4 2 4 2 4
4 1 1 2 1 3 2 4 3 1
4 1 1 2 1 3 3 4 1 2
3 2 1 2 2 3 1 3 1 3
5 3 1 2 1 3 3 4 1 5 4 2 1 2 5 4
3 2 1 2 2 3 3 1 2 1
4 2 1 2 2 3 1 4 2 3 4 3
2 1 1 2 1 2
4 1 1 2 2 3 1 4 3 4
5 3 1 2 1 3 2 4 1 5 4 1 4 2 1 2
3 2 1 2 1 3 2 3 2 1
2 1 1 2 1 2`

func parseTestcases(raw string) []testCase {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	tests := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			panic(fmt.Sprintf("line %d malformed", idx+1))
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(fmt.Sprintf("line %d bad n: %v", idx+1, err))
		}
		m, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(fmt.Sprintf("line %d bad m: %v", idx+1, err))
		}
		expectedCount := 2 + 2*(n-1+m)
		if len(parts) != expectedCount {
			panic(fmt.Sprintf("line %d expected %d numbers got %d", idx+1, expectedCount, len(parts)))
		}
		pos := 2
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			u, _ := strconv.Atoi(parts[pos])
			v, _ := strconv.Atoi(parts[pos+1])
			edges[i] = [2]int{u, v}
			pos += 2
		}
		linesPair := make([][2]int, m)
		for i := 0; i < m; i++ {
			a, _ := strconv.Atoi(parts[pos])
			b, _ := strconv.Atoi(parts[pos+1])
			linesPair[i] = [2]int{a, b}
			pos += 2
		}
		tests = append(tests, testCase{n: n, m: m, edges: edges, lines: linesPair})
	}
	if len(tests) == 0 {
		panic("no testcases parsed")
	}
	return tests
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func upd(x *pair, y pair) {
	x.l = max(x.l, y.l)
	x.r = min(x.r, y.r)
}

// Embedded solver logic from 1322F.go.
func expected(tc testCase) (int, []int) {
	n := tc.n
	repr := make([]int, n+1)
	dis := make([]int, n+1)
	val := make([]int, n+1)
	G := make([][]int, n+1)
	fa := make([][]int, n+1)
	ord := make([]int, n+1)
	dep := make([]int, n+1)
	dp := make([]int, n+1)
	typeX := make([]int, n+1)
	s := make([]pair, n+1)
	for i := 1; i <= n; i++ {
		repr[i] = i
		fa[i] = make([]int, 20)
	}
	fa[0] = make([]int, 20)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		G[u] = append(G[u], v)
		G[v] = append(G[v], u)
	}
	var find func(int) int
	find = func(x int) int {
		if repr[x] != x {
			px := repr[x]
			r := find(px)
			dis[x] ^= dis[px]
			repr[x] = r
		}
		return repr[x]
	}
	merge := func(u, v, d int) bool {
		fu := find(u)
		fv := find(v)
		if fu != fv {
			repr[fu] = fv
			dis[fu] = d ^ dis[u] ^ dis[v]
			return true
		}
		return (dis[u] ^ dis[v]) == d
	}
	kfa := func(u, k int) int {
		for i := 0; k > 0; i++ {
			if k&1 != 0 {
				u = fa[u][i]
			}
			k >>= 1
		}
		return u
	}
	var lca func(int, int) int
	lca = func(u, v int) int {
		if dep[u] < dep[v] {
			u, v = v, u
		}
		u = kfa(u, dep[u]-dep[v])
		if u == v {
			return u
		}
		for i := len(fa[u]) - 1; i >= 0; i-- {
			if fa[u][i] != fa[v][i] {
				u = fa[u][i]
				v = fa[v][i]
			}
		}
		return fa[u][0]
	}
	// DFS order and binary lifting
	stackU := []int{1}
	stackP := []int{0}
	cnt := 0
	for len(stackU) > 0 {
		u := stackU[len(stackU)-1]
		p := stackP[len(stackP)-1]
		stackU = stackU[:len(stackU)-1]
		stackP = stackP[:len(stackP)-1]
		dep[u] = dep[p] + 1
		cnt++
		ord[cnt] = u
		fa[u][0] = p
		for j := 1; j < 20; j++ {
			fa[u][j] = fa[fa[u][j-1]][j-1]
		}
		for i := len(G[u]) - 1; i >= 0; i-- {
			v := G[u][i]
			if v == p {
				continue
			}
			stackU = append(stackU, v)
			stackP = append(stackP, u)
		}
	}
	for _, line := range tc.lines {
		u, v := line[0], line[1]
		l := lca(u, v)
		if u != l {
			val[u]++
			u = kfa(u, dep[u]-dep[l]-1)
			val[u]--
		}
		if v != l {
			val[v]++
			v = kfa(v, dep[v]-dep[l]-1)
			val[v]--
		}
		if u != l && v != l {
			if !merge(u, v, 1) {
				return -1, nil
			}
		}
	}
	for i := n; i >= 1; i-- {
		u := ord[i]
		p := fa[u][0]
		val[p] += val[u]
		if val[u] != 0 {
			if !merge(u, p, 0) {
				return -1, nil
			}
		}
	}
	var solveNode func(int, int) bool
	solveNode = func(u, k int) bool {
		x := pair{1, k}
		for _, v := range G[u] {
			var t *pair
			if find(v) == find(u) {
				t = &x
			} else {
				t = &s[find(v)]
			}
			if dis[v] == dis[u] {
				typeX[v] = 0
				upd(t, pair{dp[v] + 1, k})
			} else {
				typeX[v] = 1
				upd(t, pair{1, (k + 1 - dp[v]) - 1})
			}
		}
		y := pair{1, k}
		for _, v := range G[u] {
			if find(v) != find(u) {
				comp := find(v)
				l, r := s[comp].l, s[comp].r
				if l > k+1-r {
					typeX[v] ^= 1
					nl := k + 1 - r
					nr := k + 1 - l
					s[comp].l = nl
					s[comp].r = nr
					l, r = nl, nr
				}
				upd(&y, pair{l, r})
			}
		}
		a, b := x.l, x.r
		c, d := y.l, y.r
		if max(a, c) <= min(b, d) {
			dp[u] = max(a, c)
		} else {
			c, d = k+1-d, k+1-c
			if max(a, c) <= min(b, d) {
				dp[u] = max(a, c)
				for _, v := range G[u] {
					if find(v) != find(u) {
						typeX[v] ^= 1
					}
				}
			} else {
				return false
			}
		}
		return true
	}
	check := func(k int) bool {
		for i := 1; i <= n; i++ {
			s[i].l = 1
			s[i].r = k
		}
		for i := n; i >= 1; i-- {
			if !solveNode(ord[i], k) {
				return false
			}
		}
		return true
	}
	l, r, ans := 1, n, 0
	for l <= r {
		mid := (l + r) / 2
		if check(mid) {
			ans = mid
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	check(ans)
	res := make([]int, n)
	for i := 1; i <= n; i++ {
		u := ord[i]
		typeX[u] ^= typeX[fa[u][0]]
		if typeX[u] == 0 {
			res[i-1] = dp[u]
		} else {
			res[i-1] = ans + 1 - dp[u]
		}
	}
	return ans, res
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for _, l := range tc.lines {
		sb.WriteString(fmt.Sprintf("%d %d\n", l[0], l[1]))
	}
	return sb.String()
}

func run(bin string, tc testCase) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(formatInput(tc))
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := parseTestcases(testcasesRaw)
	for idx, tc := range tests {
		wantK, wantVals := expected(tc)
		if wantK == -1 {
			fmt.Fprintf(os.Stderr, "case %d invalid constraints\n", idx+1)
			os.Exit(1)
		}
		gotStr, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(gotStr)
		if len(fields) != tc.n+1 {
			fmt.Fprintf(os.Stderr, "case %d bad output\n", idx+1)
			os.Exit(1)
		}
		gotK, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad k: %v\n", idx+1, err)
			os.Exit(1)
		}
		if gotK != wantK {
			fmt.Fprintf(os.Stderr, "case %d failed: expected k=%d got %d\n", idx+1, wantK, gotK)
			os.Exit(1)
		}
		for i := 0; i < tc.n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil || v != wantVals[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at pos %d: expected %d got %s\n", idx+1, i+1, wantVals[i], fields[1+i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
