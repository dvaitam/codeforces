package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod int = 998244353

var (
	n, k int
	adj  [][]int
)

func dfs(v, p int) []int {
	dp := make([]int, k+1)
	dp[0] = 1
	for _, to := range adj[v] {
		if to == p {
			continue
		}
		child := dfs(to, v)
		sumChild := 0
		for _, x := range child {
			sumChild += x
			if sumChild >= mod {
				sumChild -= mod
			}
		}
		newDp := make([]int, k+1)
		for i := 0; i <= k; i++ {
			if dp[i] == 0 {
				continue
			}
			val := dp[i] * sumChild % mod
			newDp[i] = (newDp[i] + val) % mod
		}
		for i := 0; i <= k; i++ {
			if dp[i] == 0 {
				continue
			}
			for j := 0; j <= k; j++ {
				if child[j] == 0 {
					continue
				}
				if i+j+1 > k {
					continue
				}
				nd := i
				if j+1 > nd {
					nd = j + 1
				}
				val := (dp[i] * child[j]) % mod
				newDp[nd] += val
				if newDp[nd] >= mod {
					newDp[nd] -= mod
				}
			}
		}
		dp = newDp
	}
	return dp
}

func solveCase(nv, kv int, edges [][2]int) string {
	n = nv
	k = kv
	adj = make([][]int, n)
	for _, e := range edges {
		u := e[0] - 1
		v := e[1] - 1
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	dp := dfs(0, -1)
	ans := 0
	for _, x := range dp {
		ans += x
		if ans >= mod {
			ans -= mod
		}
	}
	return fmt.Sprintf("%d", ans%mod)
}

type testCase struct {
	nv    int
	kv    int
	edges [][2]int
}

// Embedded testcases from testcasesF.txt.
// Format: first line is t (number of test cases).
// Each test case: first line n k, then n-1 lines each with u v (1-indexed edge).
const testcaseData = `
12
4 3
1 2
2 3
3 4
2 0
1 2
6 2
1 2
2 3
2 4
4 5
2 6
6 2
1 3
2 3
3 4
4 5
5 6
3 1
1 2
2 3
5 2
1 2
2 3
3 4
4 5
5 1
1 2
1 3
1 4
1 5
5 4
1 2
1 3
1 4
1 5
7 2
1 2
2 3
3 4
4 5
5 6
6 7
2 1
1 2
3 2
1 2
2 3
10 3
1 2
2 3
3 4
4 5
5 6
6 7
7 8
8 9
9 10
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	pos++
	res := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos+1 >= len(fields) {
			return nil, fmt.Errorf("case %d incomplete", i+1)
		}
		nv, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		kv, err := strconv.Atoi(fields[pos+1])
		if err != nil {
			return nil, fmt.Errorf("case %d bad k: %v", i+1, err)
		}
		pos += 2
		edges := make([][2]int, nv-1)
		for j := 0; j < nv-1; j++ {
			if pos+1 >= len(fields) {
				return nil, fmt.Errorf("case %d edge %d incomplete", i+1, j+1)
			}
			u, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d edge %d bad u: %v", i+1, j+1, err)
			}
			v, err := strconv.Atoi(fields[pos+1])
			if err != nil {
				return nil, fmt.Errorf("case %d edge %d bad v: %v", i+1, j+1, err)
			}
			edges[j] = [2]int{u, v}
			pos += 2
		}
		res = append(res, testCase{nv: nv, kv: kv, edges: edges})
	}
	return res, nil
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
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
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		// Build input for a single test case
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.nv, tc.kv)
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		input := sb.String()
		expected := solveCase(tc.nv, tc.kv, tc.edges)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
