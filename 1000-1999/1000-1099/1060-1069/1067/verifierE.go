package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod = int64(998244353)

// Embedded testcases previously stored in testcasesE.txt.
const testcasesEData = `100
8
1 2
1 3
3 4
4 5
4 6
2 7
1 8
3
1 2
2 3
18
1 2
1 3
1 4
3 5
3 6
2 7
7 8
2 9
5 10
4 11
1 12
11 13
13 14
5 15
13 16
9 17
7 18
6
1 2
2 3
3 4
3 5
1 6
20
1 2
2 3
3 4
2 5
2 6
2 7
4 8
5 9
2 10
9 11
5 12
1 13
5 14
10 15
12 16
10 17
17 18
7 19
14 20
14
1 2
2 3
2 4
2 5
2 6
3 7
3 8
1 9
2 10
1 11
8 12
11 13
5 14
17
1 2
2 3
1 4
2 5
1 6
4 7
2 8
8 9
5 10
3 11
6 12
7 13
12 14
10 15
6 16
7 17
11
1 2
1 3
3 4
2 5
3 6
5 7
5 8
4 9
2 10
6 11
6
1 2
2 3
1 4
1 5
3 6
3
1 2
2 3
1
11
1 2
2 3
1 4
4 5
5 6
6 7
7 8
2 9
5 10
10 11
7
1 2
2 3
1 4
3 5
4 6
5 7
6
1 2
1 3
2 4
1 5
4 6
6
1 2
2 3
2 4
1 5
4 6
7
1 2
1 3
1 4
1 5
1 6
1 7
6
1 2
1 3
3 4
4 5
5 6
8
1 2
1 3
1 4
3 5
4 6
6 7
2 8
16
1 2
1 3
2 4
4 5
4 6
1 7
2 8
7 9
8 10
4 11
11 12
7 13
4 14
8 15
4 16
2
1 2
9
1 2
1 3
3 4
2 5
2 6
4 7
7 8
5 9
5
1 2
1 3
2 4
1 5
19
1 2
1 3
2 4
4 5
1 6
4 7
2 8
3 9
6 10
5 11
11 12
8 13
13 14
11 15
6 16
14 17
17 18
7 19
9
1 2
2 3
2 4
1 5
3 6
6 7
6 8
4 9
2
1 2
20
1 2
2 3
3 4
1 5
2 6
6 7
6 8
8 9
8 10
7 11
7 12
4 13
13 14
1 15
4 16
6 17
1 18
9 19
4 20
13
1 2
1 3
3 4
1 5
2 6
2 7
6 8
6 9
9 10
8 11
9 12
8 13
1
3
1 2
1 3
16
1 2
1 3
1 4
3 5
1 6
5 7
1 8
5 9
6 10
2 11
2 12
9 13
8 14
7 15
4 16
10
1 2
1 3
2 4
4 5
1 6
1 7
1 8
6 9
9 10
14
1 2
2 3
1 4
2 5
3 6
4 7
4 8
2 9
9 10
3 11
6 12
3 13
3 14
5
1 2
2 3
2 4
3 5
18
1 2
1 3
1 4
3 5
1 6
5 7
1 8
8 9
8 10
9 11
2 12
9 13
4 14
7 15
14 16
10 17
12 18
8
1 2
1 3
3 4
1 5
5 6
3 7
5 8
15
1 2
2 3
3 4
4 5
4 6
2 7
3 8
6 9
6 10
3 11
7 12
2 13
10 14
3 15
20
1 2
2 3
2 4
2 5
5 6
3 7
6 8
2 9
2 10
7 11
11 12
3 13
6 14
11 15
6 16
11 17
6 18
10 19
1 20
20
1 2
1 3
2 4
1 5
2 6
2 7
5 8
8 9
2 10
2 11
3 12
11 13
8 14
11 15
15 16
8 17
10 18
13 19
8 20
16
1 2
2 3
2 4
2 5
3 6
5 7
6 8
8 9
7 10
9 11
7 12
6 13
5 14
8 15
7 16
19
1 2
2 3
1 4
4 5
5 6
5 7
7 8
6 9
7 10
7 11
10 12
1 13
3 14
9 15
2 16
15 17
12 18
10 19
3
1 2
2 3
8
1 2
1 3
1 4
2 5
1 6
3 7
2 8
2
1 2
13
1 2
1 3
3 4
2 5
2 6
1 7
2 8
1 9
6 10
2 11
2 12
5 13
2
1 2
20
1 2
1 3
3 4
4 5
2 6
6 7
6 8
2 9
9 10
6 11
1 12
12 13
9 14
3 15
9 16
14 17
5 18
7 19
10 20
16
1 2
2 3
1 4
2 5
3 6
5 7
4 8
5 9
7 10
6 11
3 12
7 13
1 14
7 15
1 16
9
1 2
2 3
1 4
1 5
2 6
1 7
7 8
1 9
8
1 2
1 3
2 4
2 5
4 6
6 7
6 8
13
1 2
1 3
3 4
2 5
1 6
4 7
5 8
6 9
5 10
7 11
10 12
12 13
2
1 2
9
1 2
2 3
3 4
3 5
1 6
2 7
3 8
2 9
2
1 2
17
1 2
2 3
1 4
2 5
1 6
4 7
5 8
4 9
1 10
4 11
11 12
2 13
6 14
13 15
10 16
5 17
10
1 2
1 3
1 4
2 5
5 6
3 7
5 8
1 9
1 10
1
16
1 2
1 3
2 4
2 5
3 6
3 7
7 8
7 9
7 10
1 11
3 12
7 13
8 14
2 15
13 16
17
1 2
2 3
1 4
4 5
3 6
2 7
2 8
1 9
8 10
7 11
5 12
6 13
13 14
7 15
11 16
14 17
11
1 2
2 3
2 4
1 5
1 6
4 7
3 8
5 9
3 10
10 11
11
1 2
1 3
3 4
1 5
2 6
3 7
6 8
4 9
8 10
2 11
11
1 2
2 3
3 4
2 5
2 6
1 7
2 8
2 9
2 10
6 11
10
1 2
1 3
3 4
2 5
1 6
2 7
6 8
4 9
5 10
2
1 2
11
1 2
1 3
1 4
2 5
1 6
1 7
2 8
1 9
5 10
8 11
4
1 2
1 3
2 4
16
1 2
2 3
2 4
4 5
4 6
1 7
7 8
1 9
3 10
3 11
6 12
4 13
5 14
4 15
6 16
19
1 2
2 3
2 4
3 5
5 6
6 7
1 8
4 9
4 10
6 11
3 12
9 13
4 14
1 15
1 16
4 17
7 18
10 19
11
1 2
1 3
2 4
4 5
2 6
2 7
4 8
8 9
6 10
6 11
1
15
1 2
2 3
1 4
4 5
2 6
1 7
3 8
5 9
9 10
1 11
2 12
8 13
3 14
11 15
17
1 2
1 3
3 4
2 5
5 6
3 7
5 8
1 9
1 10
7 11
7 12
4 13
13 14
10 15
3 16
2 17
5
1 2
2 3
1 4
4 5
11
1 2
1 3
3 4
3 5
5 6
6 7
5 8
6 9
5 10
10 11
11
1 2
1 3
2 4
2 5
4 6
4 7
7 8
2 9
2 10
10 11
7
1 2
1 3
3 4
2 5
5 6
1 7
18
1 2
1 3
3 4
1 5
5 6
2 7
5 8
3 9
5 10
10 11
1 12
5 13
8 14
14 15
5 16
7 17
16 18
10
1 2
1 3
2 4
4 5
1 6
1 7
3 8
6 9
2 10
7
1 2
1 3
2 4
1 5
3 6
2 7
7
1 2
1 3
3 4
1 5
1 6
1 7
1
10
1 2
1 3
1 4
3 5
2 6
6 7
2 8
1 9
1 10
3
1 2
2 3
5
1 2
2 3
3 4
3 5
9
1 2
1 3
3 4
1 5
5 6
1 7
2 8
1 9
9
1 2
2 3
3 4
3 5
2 6
5 7
3 8
8 9
8
1 2
1 3
2 4
1 5
5 6
5 7
5 8
7
1 2
2 3
2 4
2 5
3 6
2 7
11
1 2
2 3
3 4
2 5
5 6
4 7
5 8
5 9
5 10
6 11
6
1 2
1 3
1 4
3 5
2 6
6
1 2
1 3
2 4
1 5
5 6
1
2
1 2
18
1 2
1 3
1 4
3 5
4 6
4 7
2 8
8 9
7 10
5 11
9 12
2 13
6 14
13 15
1 16
4 17
6 18
2
1 2
1
13
1 2
2 3
3 4
1 5
5 6
2 7
3 8
6 9
9 10
6 11
11 12
3 13
1
12
1 2
1 3
2 4
2 5
4 6
3 7
5 8
7 9
9 10
2 11
7 12
6
1 2
2 3
2 4
2 5
1 6
3
1 2
1 3
`

type testCase struct {
	n     int
	edges [][2]int
}

// solve mirrors the logic from 1067E.go for a single test case.
func solve(n int, edges [][2]int) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	way := make([][]int64, n+1)
	ans := make([][]int64, n+1)
	for i := 1; i <= n; i++ {
		way[i] = make([]int64, 2)
		ans[i] = make([]int64, 2)
	}

	var dfs func(int, int)
	dfs = func(u, p int) {
		way[u][0] = 1
		ans[u][0] = 0
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			dfs(v, u)
			tmp1 := [2]int64{way[u][0], way[u][1]}
			tmp2 := [2]int64{ans[u][0], ans[u][1]}
			way[u][0], way[u][1], ans[u][0], ans[u][1] = 0, 0, 0, 0
			for a := 0; a < 2; a++ {
				for c := 0; c < 2; c++ {
					way[u][c] = (way[u][c] + tmp1[c]*way[v][a]) % mod
					ans[u][c] = (ans[u][c] + tmp1[c]*ans[v][a] + tmp2[c]*way[v][a]) % mod
				}
			}
			for a := 0; a < 2; a++ {
				for c := 0; c < 2; c++ {
					if a != 0 || c != 0 {
						way[u][c] = (way[u][c] + tmp1[c]*way[v][a]) % mod
						ans[u][c] = (ans[u][c] + tmp1[c]*ans[v][a] + tmp2[c]*way[v][a]) % mod
					} else {
						way[u][1] = (way[u][1] + tmp1[c]*way[v][a]) % mod
						ans[u][1] = (ans[u][1] + tmp1[c]*ans[v][a] + tmp2[c]*way[v][a] + 2*tmp1[c]*way[v][a]) % mod
					}
				}
			}
		}
	}

	dfs(1, 0)
	res := (ans[1][0] + ans[1][1]) % mod
	return fmt.Sprintf("%d", res)
}

func parseTestCases(data string) ([]testCase, error) {
	tokens := strings.Fields(data)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no embedded testcases found")
	}
	idx := 0
	t, err := strconv.Atoi(tokens[idx])
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %w", err)
	}
	idx++
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("test %d missing n", i+1)
		}
		n, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("test %d has invalid n: %w", i+1, err)
		}
		idx++
		edges := make([][2]int, n-1)
		for j := 0; j < n-1; j++ {
			if idx+1 >= len(tokens) {
				return nil, fmt.Errorf("test %d missing edge data", i+1)
			}
			u, err := strconv.Atoi(tokens[idx])
			if err != nil {
				return nil, fmt.Errorf("test %d has invalid edge u: %w", i+1, err)
			}
			v, err := strconv.Atoi(tokens[idx+1])
			if err != nil {
				return nil, fmt.Errorf("test %d has invalid edge v: %w", i+1, err)
			}
			idx += 2
			edges[j] = [2]int{u, v}
		}
		cases = append(cases, testCase{n: n, edges: edges})
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("embedded data has %d extra tokens", len(tokens)-idx)
	}
	return cases, nil
}

func runCase(bin string, tc testCase) error {
	var input strings.Builder
	input.WriteString(strconv.Itoa(tc.n))
	input.WriteByte('\n')
	for _, e := range tc.edges {
		input.WriteString(strconv.Itoa(e[0]))
		input.WriteByte(' ')
		input.WriteString(strconv.Itoa(e[1]))
		input.WriteByte('\n')
	}

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expected := solve(tc.n, tc.edges)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	cases, err := parseTestCases(testcasesEData)
	if err != nil {
		fmt.Println("failed to parse embedded testcases:", err)
		os.Exit(1)
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
