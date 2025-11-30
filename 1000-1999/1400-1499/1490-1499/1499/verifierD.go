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

type testCase struct {
	n     int
	k     int
	edges [][2]int
}

// solve mirrors 1499F.go for one test case.
func solve(tc testCase) int {
	n := tc.n
	k := tc.k
	adj := make([][]int, n)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		u--
		v--
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	var dfs func(v, p int) []int
	dfs = func(v, p int) []int {
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
			// cut edge v-to
			for i := 0; i <= k; i++ {
				if dp[i] == 0 {
					continue
				}
				val := dp[i] * sumChild % mod
				newDp[i] = (newDp[i] + val) % mod
			}
			// keep edge v-to
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

	dp := dfs(0, -1)
	ans := 0
	for _, x := range dp {
		ans += x
		if ans >= mod {
			ans -= mod
		}
	}
	return ans % mod
}

// Embedded testcases from testcasesF.txt.
const testcaseData = `
100
6 3
1 2
1 3
2 4
2 5
1 6
3 1
1 2
2 3
3 2
1 2
1 3
2 1
1 2
4 2
1 2
1 3
1 4
3 3
1 2
1 3
3 1
1 2
1 3
3 1
1 2
2 3
4 2
1 2
1 3
3 4
3 2
1 2
1 3
4 4
1 2
1 3
2 4
2 2
1 2
6 5
1 2
2 3
1 4
3 5
3 6
4 4
1 2
1 3
2 4
5 2
1 2
2 3
1 4
3 5
5 1
1 2
2 3
2 4
1 5
5 1
1 2
1 3
1 4
2 5
5 3
1 2
2 3
2 4
1 5
6 6
1 2
2 3
1 4
4 5
1 6
3 2
1 2
1 3
4 3
1 2
2 3
3 4
4 3
1 2
1 3
3 4
3 3
1 2
2 3
3 3
1 2
1 3
6 5
1 2
1 3
1 4
3 5
3 6
5 4
1 2
1 3
3 4
1 5
5 3
1 2
1 3
3 4
2 5
5 1
1 2
2 3
2 4
2 5
2 2
1 2
3 2
1 2
2 3
5 3
1 2
2 3
2 4
4 5
5 2
1 2
2 3
3 4
2 5
5 3
1 2
1 3
2 4
3 5
6 5
1 2
1 3
2 4
1 5
3 6
4 3
1 2
2 3
3 4
4 3
1 2
1 3
2 4
6 3
1 2
2 3
2 4
3 5
5 6
4 3
1 2
2 3
3 4
5 3
1 2
2 3
2 4
3 5
6 2
1 2
1 3
2 4
4 5
3 6
2 2
1 2
6 5
1 2
2 3
3 4
3 5
1 6
3 1
1 2
1 3
3 3
1 2
1 3
5 2
1 2
1 3
1 4
1 5
4 2
1 2
1 3
3 4
2 2
1 2
4 4
1 2
1 3
1 4
4 1
1 2
2 3
2 4
5 1
1 2
1 3
3 4
3 5
2 2
1 2
6 3
1 2
1 3
2 4
4 5
5 6
6 3
1 2
1 3
1 4
4 5
5 6
5 2
1 2
1 3
2 4
1 5
5 1
1 2
1 3
2 4
4 5
5 4
1 2
2 3
2 4
4 5
5 4
1 2
2 3
3 4
3 5
5 3
1 2
1 3
2 4
3 5
5 3
1 2
2 3
1 4
1 5
6 4
1 2
2 3
1 4
2 5
4 6
2 2
1 2
3 2
1 2
1 3
5 3
1 2
2 3
2 4
4 5
6 4
1 2
1 3
1 4
3 5
3 6
6 6
1 2
2 3
3 4
1 5
4 6
4 3
1 2
1 3
1 4
4 4
1 2
1 3
1 4
5 4
1 2
2 3
2 4
4 5
3 2
1 2
1 3
4 1
1 2
2 3
3 4
5 3
1 2
1 3
2 4
4 5
4 2
1 2
1 3
1 4
3 2
1 2
1 3
2 2
1 2
2 2
1 2
2 1
1 2
3 3
1 2
1 3
4 4
1 2
1 3
3 4
3 2
1 2
1 3
6 2
1 2
1 3
3 4
4 5
1 6
4 1
1 2
2 3
3 4
3 1
1 2
1 3
4 4
1 2
2 3
3 4
4 1
1 2
2 3
1 4
6 6
1 2
1 3
2 4
1 5
5 6
3 2
1 2
1 3
3 2
1 2
2 3
3 1
1 2
2 3
2 1
1 2
2 2
1 2
5 4
1 2
2 3
1 4
1 5
6 6
1 2
1 3
1 4
4 5
2 6
2 1
1 2
3 1
1 2
2 3
3 2
1 2
2 3
5 4
1 2
2 3
3 4
3 5
5 3
1 2
1 3
3 4
2 5
5 1
1 2
1 3
1 4
2 5
3 1
1 2
2 3
5 3
1 2
2 3
2 4
3 5
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
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if pos+1 >= len(fields) {
			return nil, fmt.Errorf("case %d missing n/k", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", caseIdx+1, err)
		}
		pos++
		k, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d bad k: %v", caseIdx+1, err)
		}
		pos++
		if pos+2*(n-1) > len(fields) {
			return nil, fmt.Errorf("case %d missing edges", caseIdx+1)
		}
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			u, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d edge %d bad u: %v", caseIdx+1, i+1, err)
			}
			v, err := strconv.Atoi(fields[pos+1])
			if err != nil {
				return nil, fmt.Errorf("case %d edge %d bad v: %v", caseIdx+1, i+1, err)
			}
			edges[i] = [2]int{u, v}
			pos += 2
		}
		res = append(res, testCase{n: n, k: k, edges: edges})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra tokens at end")
	}
	return res, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	input := sb.String()

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
		expected := strconv.Itoa(solve(tc))
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, fmt.Sprintf("%d %d ...edges omitted...", tc.n, tc.k), expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
