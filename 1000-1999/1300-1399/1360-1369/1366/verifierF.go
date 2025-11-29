package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const modVal int64 = 1_000_000_007
const infNeg int64 = -1 << 60

type edge struct {
	u, v int
	w    int64
}

type line struct {
	k int64
	b int64
}

type testcase struct {
	n, m, q int
	edges   []edge
}

// Embedded copy of testcasesF.txt so the verifier is self-contained.
const testcasesRaw = `4 5 10 3 1 4 4 1 1 1 2 3 2 4 4 2 1 5
4 3 7 2 1 1 3 1 4 4 2 2
3 2 2 3 2 2 2 1 5
2 1 1 2 1 2
2 1 3 2 1 2
4 4 5 3 2 1 4 2 4 1 3 3 2 1 4
3 2 6 2 1 1 3 2 3
3 3 8 1 3 1 2 1 2 3 2 4
3 2 7 3 1 5 2 1 4
3 3 6 1 3 1 3 1 2 2 1 4
2 1 3 2 1 5
3 3 3 3 1 4 2 1 3 2 3 5
3 2 4 3 1 3 2 1 5
4 5 7 1 2 1 2 1 1 3 1 3 1 4 2 4 2 2
3 3 4 3 1 4 1 2 5 2 1 5
2 1 4 2 1 4
3 3 6 3 1 4 2 1 2 1 3 4
3 2 2 3 2 1 2 1 3
3 2 4 3 2 4 2 1 5
3 3 5 3 1 4 2 3 3 2 1 5
4 5 8 3 4 3 3 2 3 3 4 2 2 1 5 4 2 3
4 4 8 2 4 5 2 1 2 3 2 4 4 2 1
4 6 8 1 2 1 2 1 1 4 3 4 2 3 4 3 1 2 1 2 4
2 1 4 2 1 3
3 2 6 3 2 1 2 1 2
3 3 5 2 1 1 1 2 3 3 2 1
4 5 9 3 1 4 2 4 5 4 2 5 2 1 3 3 4 2
3 2 3 3 1 4 2 1 3
2 1 5 2 1 3
4 6 9 2 4 5 3 2 3 4 2 4 2 1 2 1 3 3 3 4 4
3 3 4 2 1 1 3 1 5 1 2 3
2 1 4 2 1 4
4 5 6 2 1 1 3 1 4 2 3 3 1 4 3 4 2 2
3 2 2 3 2 1 2 1 3
4 4 4 2 3 4 4 2 4 3 2 4 2 1 4
2 1 1 2 1 5
3 3 5 1 3 4 2 1 2 3 2 5
4 3 3 4 3 1 2 1 3 3 1 1
3 3 3 1 3 2 3 1 1 2 1 4
3 3 4 3 1 4 1 3 1 2 1 5
4 3 8 4 3 5 2 1 1 3 2 1
2 1 5 2 1 1
3 2 4 3 2 3 2 1 4
2 1 6 2 1 2
4 5 9 2 3 2 2 1 3 4 1 3 3 1 5 3 4 4
2 1 6 2 1 3
2 1 1 2 1 5
4 6 7 4 1 5 3 2 3 3 4 2 2 1 4 2 4 2 1 4 1
3 3 7 3 2 2 2 3 3 2 1 4
4 6 8 4 2 1 1 2 2 3 1 2 2 1 5 3 4 1 1 2 4
3 3 5 1 3 2 2 1 3 3 2 4
2 1 4 2 1 2
4 5 5 4 1 4 3 2 2 1 3 1 2 1 2 1 3 3
4 6 9 1 2 1 2 1 1 3 4 3 3 2 3 1 2 2 4 2 2
2 1 3 2 1 1
4 6 7 1 4 4 4 3 5 2 4 1 2 1 2 3 4 1 3 1 1
2 1 1 2 1 2
3 2 5 3 1 4 2 1 2
3 2 2 3 2 3 2 1 4
4 3 3 4 3 2 3 2 2 2 1 2
2 1 1 2 1 5
2 1 6 2 1 5
2 1 5 2 1 4
3 3 7 2 1 1 1 3 1 3 2 4
2 1 4 2 1 3
3 3 4 3 1 4 1 3 1 2 1 5
3 2 4 3 1 1 2 1 4
4 4 6 2 1 1 4 2 4 1 3 2 3 1 4
4 3 6 2 1 1 3 1 4 4 2 1
3 2 4 2 1 3 3 1 2
3 3 8 1 3 1 2 1 3 3 1 1
2 1 2 2 1 4
3 2 3 3 2 2 2 1 2
3 3 3 3 1 1 2 3 5 2 1 5
2 1 4 2 1 4
2 1 3 2 1 2
3 3 4 2 1 1 3 1 5 1 3 1
2 1 6 2 1 2
2 1 5 2 1 4
4 5 8 2 3 3 3 1 3 2 1 2 4 1 3 2 3 1
4 6 9 2 3 3 4 2 4 2 4 1 2 1 2 3 1 2 1 2 4
4 3 3 3 1 3 4 3 1 2 1 4
4 5 9 4 3 4 3 2 5 2 1 5 2 4 3 1 4 5
4 3 8 3 2 1 4 1 1 2 1 3
3 3 5 2 3 3 2 1 3 3 1 1
4 4 7 3 1 3 2 1 1 1 4 4 4 3 4
3 2 7 3 2 4 2 1 5
3 2 7 2 1 2 3 1 2
3 3 5 3 1 4 2 1 2 1 3 5
4 6 9 3 4 3 4 3 2 1 3 2 3 2 5 1 4 2 2 1 4
2 1 5 2 1 4
2 1 5 2 1 3
4 3 3 3 1 3 2 1 2 4 1 4
4 6 7 2 1 1 4 3 2 2 3 3 1 2 3 1 4 3 3 1 2
2 1 2 2 1 3
3 2 5 2 1 2 3 2 5
3 3 3 3 2 2 1 3 4 2 1 5
3 3 8 3 2 1 1 2 3 2 1 4
4 4 6 4 3 1 2 1 2 3 2 4 1 3 4
3 2 6 2 1 1 3 1 4`

func parseTestcases() ([]testcase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testcase, 0, len(lines))
	for idx, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		fields := strings.Fields(ln)
		if len(fields) < 3 {
			return nil, fmt.Errorf("line %d: expected at least 3 fields", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %v", idx+1, err)
		}
		q, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse q: %v", idx+1, err)
		}
		if (len(fields)-3)%3 != 0 || len(fields) != 3+3*m {
			return nil, fmt.Errorf("line %d: expected %d fields, got %d", idx+1, 3+3*m, len(fields))
		}
		edges := make([]edge, m)
		for i := 0; i < m; i++ {
			u, err := strconv.Atoi(fields[3+3*i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse u%d: %v", idx+1, i+1, err)
			}
			v, err := strconv.Atoi(fields[4+3*i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse v%d: %v", idx+1, i+1, err)
			}
			w, err := strconv.ParseInt(fields[5+3*i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse w%d: %v", idx+1, i+1, err)
			}
			edges[i] = edge{u: u, v: v, w: w}
		}
		cases = append(cases, testcase{n: n, m: m, q: q, edges: edges})
	}
	return cases, nil
}

func isBad(a, b, c line) bool {
	leftNum := (b.b - a.b) * (b.k - c.k)
	rightNum := (c.b - b.b) * (a.k - b.k)
	return leftNum >= rightNum
}

func floorDiv(a, b int64) int64 {
	if b < 0 {
		a, b = -a, -b
	}
	if a >= 0 {
		return a / b
	}
	return -((-a + b - 1) / b)
}

func mod(x int64) int64 {
	x %= modVal
	if x < 0 {
		x += modVal
	}
	return x
}

// solve implements the logic from 1366F.go for a single testcase.
func solve(tc testcase) string {
	n, m, q := tc.n, tc.m, tc.q
	edges := tc.edges

	maxW := make([]int64, n+1)
	for _, e := range edges {
		if e.w > maxW[e.u] {
			maxW[e.u] = e.w
		}
		if e.w > maxW[e.v] {
			maxW[e.v] = e.w
		}
	}

	if q <= m {
		dp := make([][]int64, q+1)
		for i := range dp {
			dp[i] = make([]int64, n+1)
			for j := 1; j <= n; j++ {
				dp[i][j] = infNeg
			}
		}
		dp[0][1] = 0
		ans := int64(0)
		for t := 1; t <= q; t++ {
			for j := 1; j <= n; j++ {
				dp[t][j] = infNeg
			}
			for _, e := range edges {
				if dp[t-1][e.u] != infNeg {
					if val := dp[t-1][e.u] + e.w; val > dp[t][e.v] {
						dp[t][e.v] = val
					}
				}
				if dp[t-1][e.v] != infNeg {
					if val := dp[t-1][e.v] + e.w; val > dp[t][e.u] {
						dp[t][e.u] = val
					}
				}
			}
			best := infNeg
			for i := 1; i <= n; i++ {
				if dp[t][i] > best {
					best = dp[t][i]
				}
			}
			ans = (ans + mod(best)) % modVal
		}
		return strconv.FormatInt(ans, 10)
	}

	// q > m
	dp := make([][]int64, m+1)
	for i := range dp {
		dp[i] = make([]int64, n+1)
		for j := 1; j <= n; j++ {
			dp[i][j] = infNeg
		}
	}
	dp[0][1] = 0
	partial := int64(0)
	for t := 1; t <= m; t++ {
		for j := 1; j <= n; j++ {
			dp[t][j] = infNeg
		}
		for _, e := range edges {
			if dp[t-1][e.u] != infNeg {
				if val := dp[t-1][e.u] + e.w; val > dp[t][e.v] {
					dp[t][e.v] = val
				}
			}
			if dp[t-1][e.v] != infNeg {
				if val := dp[t-1][e.v] + e.w; val > dp[t][e.u] {
					dp[t][e.u] = val
				}
			}
		}
		best := infNeg
		for i := 1; i <= n; i++ {
			if dp[t][i] > best {
				best = dp[t][i]
			}
		}
		partial = (partial + mod(best)) % modVal
	}

	lineMap := make(map[int64]int64)
	for v := 1; v <= n; v++ {
		slope := maxW[v]
		if slope == 0 {
			continue
		}
		best := infNeg
		for t := 0; t <= m; t++ {
			if dp[t][v] == infNeg {
				continue
			}
			val := dp[t][v] - int64(t)*slope
			if val > best {
				best = val
			}
		}
		if cur, ok := lineMap[slope]; !ok || best > cur {
			lineMap[slope] = best
		}
	}

	lines := make([]line, 0, len(lineMap))
	for k, b := range lineMap {
		lines = append(lines, line{k: int64(k), b: b})
	}
	sort.Slice(lines, func(i, j int) bool {
		if lines[i].k == lines[j].k {
			return lines[i].b > lines[j].b
		}
		return lines[i].k < lines[j].k
	})
	hull := make([]line, 0)
	for _, ln := range lines {
		if len(hull) > 0 && hull[len(hull)-1].k == ln.k {
			if hull[len(hull)-1].b >= ln.b {
				continue
			}
			hull[len(hull)-1] = ln
			continue
		}
		for len(hull) >= 2 && isBad(hull[len(hull)-2], hull[len(hull)-1], ln) {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, ln)
	}

	pos := int64(m + 1)
	ans := partial
	for i := 0; i < len(hull) && pos <= int64(q); i++ {
		right := int64(q)
		if i+1 < len(hull) {
			num := hull[i].b - hull[i+1].b
			den := hull[i+1].k - hull[i].k
			cross := floorDiv(num, den)
			if cross < right {
				right = cross
			}
		}
		if right < pos {
			continue
		}
		if right > int64(q) {
			right = int64(q)
		}
		count := right - pos + 1
		sumX := (pos + right) * count / 2
		term := (mod(hull[i].k) * mod(sumX)) % modVal
		term = (term + mod(hull[i].b)*mod(count)) % modVal
		ans = (ans + term) % modVal
		pos = right + 1
	}
	return strconv.FormatInt(ans, 10)
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.q)
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
		}
		input := sb.String()
		expect := solve(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
