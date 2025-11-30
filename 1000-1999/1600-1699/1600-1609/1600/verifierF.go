package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesF.txt so the verifier is self-contained.
const testcasesRaw = `5 9 2 4 1 2 1 5 1 4 2 3 4 5 2 5 1 3 3 5
6 6 2 4 1 5 4 6 2 3 4 5 5 6
7 11 2 4 3 4 4 6 1 4 5 7 1 7 5 6 3 6 1 6 2 5 4 7
7 11 3 4 2 7 5 7 1 4 4 5 5 6 3 6 1 6 2 5 4 7 3 5
6 11 1 2 3 4 1 5 1 4 2 3 4 5 2 6 3 6 1 6 1 3 3 5
7 9 1 2 3 4 4 6 1 4 2 3 5 6 2 5 1 3 4 7
5 6 2 4 1 2 1 4 4 5 2 5 3 5
7 20 3 4 3 7 5 7 1 6 2 5 1 3 4 5 5 6 3 6 2 4 1 2 2 7 1 5 6 7 4 7 3 5 1 4 2 3 1 7 2 6
5 1 3 5
7 8 3 4 3 7 4 6 1 4 4 5 2 6 5 6 1 3
5 4 2 3 4 5 2 5 1 3
6 2 5 6 1 4
7 10 2 4 1 2 1 5 3 7 5 7 2 3 1 7 2 6 1 6 1 3
5 2 2 3 1 3
7 16 2 4 1 2 3 7 5 7 1 4 2 3 6 7 1 7 4 5 2 6 5 6 3 6 1 6 2 5 4 7 3 5
7 2 2 4 1 2
7 10 1 2 3 4 2 7 4 6 5 7 2 3 1 7 2 6 5 6 3 6
6 9 3 4 1 5 4 6 4 5 2 6 3 6 1 6 2 5 3 5
6 3 1 2 5 6 3 6
7 2 4 7 2 7
5 7 1 2 3 4 2 7 4 6 5 6 3 6 2 5
5 2 2 3 1 3
7 1 1 4
6 5 3 4 2 3 2 6 1 2 1 4
7 11 3 4 1 2 2 7 5 7 2 3 1 3 2 6 1 6 1 5 4 7 2 5
7 14 3 4 2 7 5 7 1 6 1 2 2 6 5 6 4 5 1 3 2 5 2 3 4 7 3 5 1 4
5 5 1 2 3 4 1 5 4 5 3 4 4 5
6 13 2 4 1 2 1 5 1 4 2 3 4 5 5 6 1 6 4 6 3 6 2 5 3 4 3 5
7 9 1 2 3 4 2 7 4 6 4 5 2 6 3 6 1 3 4 7
7 12 1 2 3 4 1 3 4 6 4 5 5 6 1 6 2 5 3 4 1 5 2 3 3 5
6 12 1 2 3 4 1 5 1 4 2 3 4 6 1 6 2 6 5 6 2 5 3 4 3 5
7 11 2 4 1 2 1 5 3 7 5 7 3 4 2 3 2 6 5 6 4 7 2 5 3 5
6 3 1 2 4 6 2 6
6 9 1 2 3 4 1 5 1 4 2 3 1 6 3 6 4 6 2 5 3 5
7 9 1 2 3 4 1 3 5 7 5 7 4 5 3 6 2 4 1 6
7 10 1 2 3 4 1 5 1 3 4 6 4 5 5 6 3 6 1 6 4 7 2 5
7 11 1 2 3 4 1 3 1 5 4 6 4 5 1 6 2 6 2 5 3 4 2 3 3 5
6 10 1 2 1 3 2 4 1 4 2 3 1 6 3 4 4 6 3 6 2 5 3 5
6 3 3 4 3 6 1 6
5 6 2 4 3 4 4 5 3 5 1 3 2 5
6 9 1 2 2 7 5 7 2 3 1 3 3 4 1 5 2 6 3 6 4 5
5 5 1 2 1 3 1 4 1 5 3 5 4 5
6 8 3 4 1 2 4 5 1 5 1 4 2 3 3 6 1 6 3 5
5 4 1 2 1 3 2 4 1 5 3 5 2 5
6 7 1 2 2 7 1 7 4 6 4 5 3 6 1 6 1 5
6 3 2 4 1 2 2 5
6 6 1 2 4 7 2 7 4 5 5 6 1 6 3 5
5 2 4 7 2 7
6 4 1 2 1 3 2 3 5 6 3 6
5 7 1 2 1 5 1 4 2 6 2 5 4 5
6 5 1 2 3 4 1 5 4 6 3 6 1 6
6 4 1 2 1 5 3 4 3 6 4 6
5 4 1 2 1 3 2 4 2 5 4 5 3 5
5 7 1 2 1 5 3 4 2 3 4 5 3 5 4 5
6 4 1 2 2 3 4 5 3 5 2 6 3 6
5 7 1 2 1 3 3 4 2 4 2 6 3 5 2 5
7 10 1 2 3 4 4 6 2 7 1 3 1 7 4 5 2 6 3 4 2 5 3 6
7 12 1 2 3 4 1 3 2 7 4 5 5 6 1 5 3 6 2 6 3 4 2 5 4 7 1 7
7 8 1 2 3 4 2 7 4 6 2 5 3 6 1 6 1 3 1 4
6 3 1 2 2 5 4 5
6 6 1 2 2 3 4 5 2 5 1 5 1 6
7 11 1 2 3 4 1 5 3 4 4 7 5 6 2 6 1 6 1 3 2 5 3 5 2 3
7 10 1 2 1 5 1 3 4 6 4 7 5 7 1 4 1 6 2 5 3 5 2 3
5 9 1 2 3 4 1 4 1 5 4 5 3 5 2 5 2 3 1 3
7 11 2 4 1 2 1 3 3 7 2 7 4 6 4 5 5 6 2 3 2 6 3 5 3 6
7 12 1 2 3 4 2 7 1 5 4 5 3 6 2 4 1 6 2 6 4 7 3 5 2 5
7 12 1 2 1 3 1 4 4 7 2 6 3 6 5 6 2 5 3 5 2 3 2 4 1 5 4 5
7 12 1 2 1 3 1 4 5 7 1 6 4 6 3 5 4 5 1 5 2 6 3 6 2 4 2 5
7 12 1 2 2 7 4 7 1 5 1 4 1 3 1 6 4 6 3 6 2 6 2 5 4 5 1 4
7 8 3 4 4 6 1 4 1 3 2 5 1 6 3 6 2 3
7 11 1 2 1 5 1 6 1 3 4 5 2 4 4 7 3 6 2 6 2 5 4 6 3 5
7 7 1 2 1 3 1 4 2 6 4 6 1 5 3 5 2 5
7 11 3 4 2 7 1 4 4 6 5 7 3 6 2 6 2 5 3 5 1 3 1 6 3 4
7 7 2 7 4 6 2 6 4 5 2 5 1 4 1 3
7 7 3 4 2 7 5 7 3 6 1 5 2 4 3 4
6 5 1 2 3 4 3 4 2 6 3 6 1 6
6 4 1 2 2 3 1 4 1 3 3 6
6 4 1 2 1 3 2 4 2 5 3 5
6 7 1 2 3 4 2 3 4 5 1 5 3 5 1 6
5 3 1 2 1 3 2 4
5 4 3 4 2 3 1 4 2 5 3 5
5 3 2 4 1 4 1 5 2 5
6 5 1 2 2 7 4 6 2 6 3 6 1 6
5 2 3 4 1 5
7 6 1 2 3 4 1 4 3 6 2 5 1 5 3 5 2 3
7 10 1 2 1 5 3 4 4 6 5 7 1 4 2 4 3 6 2 6 2 5 1 3
7 10 2 7 4 6 1 6 1 4 5 7 5 6 3 4 2 5 2 3 1 3
7 6 2 4 3 4 1 5 5 6 1 3 3 5
6 6 2 4 1 2 1 5 1 4 2 3 3 6 2 6 3 5
6 8 1 2 1 5 3 4 2 3 4 6 2 6 3 6 1 6
6 6 1 2 2 3 1 4 3 4 2 6 3 6
6 7 1 2 3 4 2 3 4 5 3 5 1 5 2 5
6 8 1 2 1 3 1 4 4 5 1 6 2 6 3 6 2 5
6 4 1 2 1 3 2 4 2 6 2 5
6 7 1 2 3 4 1 4 1 6 3 6 2 6 3 5
6 5 1 2 1 3 2 4 3 5 1 5
6 4 2 7 1 7 4 5 1 3 2 6
6 4 1 2 2 3 3 4 1 6 2 6
6 2 1 2 1 3
6 7 1 2 1 3 2 3 3 4 1 5 3 5 2 6
6 5 1 2 1 3 1 5 2 6 4 6 3 5
6 6 1 2 3 4 2 4 2 6 3 6 2 5
6 5 1 2 1 3 1 4 1 5 3 6
6 2 3 4 1 3
6 6 1 2 1 3 2 3 2 6 1 5 2 4
6 6 1 2 1 3 2 4 1 4 2 6 3 6
6 6 1 2 1 3 2 3 1 5 2 6 3 6
6 5 1 2 3 4 2 4 3 5 2 5 1 5
6 8 1 2 1 3 2 4 2 6 1 5 3 5 4 5 3 6
6 4 1 2 2 3 1 3 2 6 3 5
6 4 1 2 1 3 2 3 2 4 1 5 3 5
5 2 1 2 1 3
6 6 1 2 1 3 1 5 2 6 4 6 2 5
5 2 2 3 1 2
6 4 1 2 1 5 2 6 3 6
6 4 1 2 1 4 2 3 3 6 2 6
6 4 1 2 1 3 2 4 2 6 3 5
6 4 1 2 1 3 1 4 1 5 3 5
6 6 1 2 1 3 1 4 1 6 2 6 3 6
6 5 1 2 1 3 1 4 1 5 3 4
6 3 1 2 1 4 1 5
6 4 1 2 1 3 1 5 2 5 3 5
6 5 1 2 1 3 1 4 1 5 2 6
6 5 1 2 1 4 1 5 2 5 3 5
6 4 1 2 1 3 1 5 2 5 3 6
6 5 1 2 1 5 3 4 1 3 2 5
6 2 2 4 1 3
6 3 1 2 1 4 1 3
6 5 1 2 1 3 2 4 2 5 3 4
6 3 1 2 1 3 2 5
6 5 1 2 1 3 2 5 3 4 1 4
6 4 2 4 1 2 1 3 2 5 3 4
6 3 3 4 2 6 1 6
6 3 1 2 2 4 2 3
6 4 1 2 2 3 1 4 2 4 2 6
6 5 1 2 1 3 2 4 1 4 3 6
6 5 1 2 1 3 1 4 2 6 3 6
6 3 1 2 2 4 3 5
6 4 1 2 1 3 2 4 1 4 2 5
6 3 2 3 1 3
6 3 1 2 2 4 3 6
5 3 1 2 1 3 1 4`

type testCase struct {
	n     int
	edges [][2]int
}

// solveCase replicates the 1600F logic, searching for a 5-clique or 5-independent set.
func solveCase(tc testCase) ([]int, bool) {
	if tc.n < 5 {
		return nil, false
	}
	limit := tc.n
	if limit > 43 {
		limit = 43
	}
	adj := make([][]bool, limit)
	for i := 0; i < limit; i++ {
		adj[i] = make([]bool, limit)
	}
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		if u <= limit && v <= limit {
			adj[u-1][v-1] = true
			adj[v-1][u-1] = true
		}
	}
	ids := make([]int, limit)
	for i := 0; i < limit; i++ {
		ids[i] = i + 1
	}
	choose := make([]int, 5)
	var ans []int
	var dfs func(start, depth int)
	dfs = func(start, depth int) {
		if len(ans) > 0 {
			return
		}
		if depth == 5 {
			clique := true
			independent := true
			for i := 0; i < 5; i++ {
				for j := i + 1; j < 5; j++ {
					if adj[choose[i]-1][choose[j]-1] {
						independent = false
					} else {
						clique = false
					}
				}
			}
			if clique || independent {
				ans = append([]int(nil), choose...)
			}
			return
		}
		for i := start; i <= limit-(5-depth); i++ {
			choose[depth] = ids[i]
			dfs(i+1, depth+1)
			if len(ans) > 0 {
				return
			}
		}
	}
	dfs(0, 0)
	if len(ans) == 0 {
		return nil, false
	}
	return ans, true
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %v", idx+1, err)
		}
		rest := fields[2:]
		if len(rest)%2 != 0 {
			rest = rest[:len(rest)-1]
		}
		actual := len(rest) / 2
		if actual < m {
			m = actual
		}
		edges := make([][2]int, m)
		pos := 0
		for i := 0; i < m; i++ {
			u, err := strconv.Atoi(rest[pos])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse u: %v", idx+1, err)
			}
			v, err := strconv.Atoi(rest[pos+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse v: %v", idx+1, err)
			}
			edges[i] = [2]int{u, v}
			pos += 2
		}
		cases = append(cases, testCase{n: n, edges: edges})
	}
	return cases, nil
}

func buildEdgeMap(tc testCase) map[int]map[int]bool {
	m := make(map[int]map[int]bool)
	for _, e := range tc.edges {
		a, b := e[0], e[1]
		if m[a] == nil {
			m[a] = make(map[int]bool)
		}
		if m[b] == nil {
			m[b] = make(map[int]bool)
		}
		m[a][b] = true
		m[b][a] = true
	}
	return m
}

func verifyOutput(out string, tc testCase, feasible bool, full map[int]map[int]bool) error {
	out = strings.TrimSpace(out)
	if out == "-1" {
		if feasible {
			return fmt.Errorf("expected subset but got -1")
		}
		return nil
	}
	fields := strings.Fields(out)
	if len(fields) != 5 {
		return fmt.Errorf("expected 5 numbers got %d", len(fields))
	}
	nums := make([]int, 5)
	seen := make(map[int]bool)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid number: %v", err)
		}
		if v < 1 || v > tc.n {
			return fmt.Errorf("number out of range")
		}
		if seen[v] {
			return fmt.Errorf("duplicate numbers")
		}
		seen[v] = true
		nums[i] = v
	}
	clique := true
	independent := true
	for i := 0; i < 5; i++ {
		for j := i + 1; j < 5; j++ {
			if full[nums[i]][nums[j]] {
				independent = false
			} else {
				clique = false
			}
		}
	}
	if !(clique || independent) {
		return fmt.Errorf("numbers do not form clique or independent set")
	}
	if !feasible {
		return fmt.Errorf("output subset but none exist")
	}
	return nil
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

	for i, tc := range cases {
		expectedSet, feasible := solveCase(tc)
		fullEdges := buildEdgeMap(tc)

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verifyOutput(got, tc, feasible, fullEdges); err != nil {
			// if candidate output matches exact expected set, also accept
			if feasible && strings.TrimSpace(got) == strings.TrimSpace(strings.Trim(fmt.Sprint(expectedSet), "[]")) {
				continue
			}
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
