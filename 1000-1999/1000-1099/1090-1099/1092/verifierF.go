package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
10 2 8 5 1 1 3 10 8 6 6 1 2 2 3 2 4 2 5 4 6 5 7 5 8 2 9 4 10
10 9 5 10 2 7 6 2 6 7 5 1 2 1 3 1 4 3 5 1 6 1 7 5 8 4 9 6 10
8 4 9 10 9 1 6 4 10 1 2 2 3 2 4 1 5 1 6 5 7 6 8
9 4 2 10 5 5 4 7 8 4 1 2 1 3 3 4 1 5 2 6 2 7 1 8 6 9
9 10 10 5 6 7 9 7 5 3 1 2 1 3 1 4 4 5 5 6 6 7 4 8 2 9
8 4 2 10 8 8 7 2 9 1 2 2 3 2 4 4 5 1 6 2 7 6 8
5 8 8 3 1 1 1 2 2 3 3 4 3 5
4 5 9 8 6 1 2 2 3 2 4
10 8 5 10 8 8 8 3 7 8 5 1 2 2 3 2 4 2 5 5 6 4 7 6 8 5 9 6 10
7 8 3 5 9 1 10 8 1 2 1 3 1 4 1 5 3 6 3 7
8 10 1 4 3 5 10 1 7 1 2 1 3 2 4 2 5 1 6 3 7 3 8
3 4 10 6 1 2 1 3
2 2 1 1 2
8 2 1 9 3 10 2 3 10 1 2 2 3 2 4 1 5 2 6 2 7 4 8
8 8 8 9 5 8 10 2 5 1 2 2 3 2 4 1 5 1 6 5 7 7 8
5 5 5 8 7 8 1 2 1 3 3 4 1 5
4 10 4 9 1 1 2 2 3 2 4
1 1
4 7 6 2 3 1 2 2 3 3 4
3 8 8 2 1 2 1 3
6 9 1 3 9 3 3 1 2 1 3 3 4 4 5 2 6
4 10 3 7 3 1 2 1 3 3 4
10 1 10 7 4 6 7 9 8 1 7 1 2 1 3 3 4 3 5 3 6 6 7 2 8 3 9 3 10
3 7 2 1 1 2 1 3
8 1 2 9 2 6 10 1 1 1 2 2 3 1 4 3 5 1 6 4 7 1 8
10 1 3 8 2 2 3 1 1 7 5 1 2 2 3 1 4 2 5 5 6 6 7 2 8 5 9 1 10
5 9 6 9 2 10 1 2 2 3 2 4 1 5
4 1 6 6 7 1 2 2 3 1 4
1 1
6 2 8 1 8 5 2 1 2 1 3 1 4 2 5 4 6
9 3 4 7 10 7 10 9 8 5 1 2 1 3 2 4 1 5 2 6 4 7 1 8 3 9
10 2 1 9 3 2 3 10 7 8 8 1 2 1 3 3 4 4 5 2 6 4 7 7 8 1 9 6 10
7 5 6 3 5 6 10 6 1 2 1 3 2 4 2 5 1 6 4 7
10 10 4 2 4 7 7 3 3 9 1 1 2 2 3 2 4 4 5 3 6 6 7 6 8 2 9 5 10
3 1 2 10 1 2 1 3
2 2 2 1 2
6 6 3 4 3 3 4 1 2 2 3 3 4 4 5 5 6
2 10 3 1 2
4 10 6 3 6 1 2 1 3 2 4
1 1
1 8
8 6 1 1 5 2 7 10 6 1 2 2 3 1 4 2 5 3 6 6 7 7 8
8 7 1 9 3 2 8 6 1 1 2 2 3 2 4 1 5 4 6 1 7 4 8
4 9 3 8 6 1 2 2 3 2 4
9 5 8 9 1 9 4 10 1 9 1 2 2 3 1 4 1 5 3 6 1 7 7 8 4 9
10 5 4 3 5 1 5 5 6 4 3 1 2 1 3 2 4 2 5 3 6 2 7 4 8 4 9 8 10
7 9 10 6 2 2 7 1 1 2 2 3 2 4 2 5 1 6 4 7
9 9 5 9 5 6 5 7 6 2 1 2 1 3 2 4 1 5 5 6 3 7 7 8 1 9
2 2 1 1 2
7 9 2 3 5 6 1 2 1 2 2 3 2 4 4 5 4 6 5 7
2 2 5 1 2
8 5 5 4 1 1 9 9 7 1 2 1 3 2 4 1 5 1 6 3 7 7 8
3 9 2 1 1 2 2 3
8 4 2 5 3 3 3 2 1 1 2 2 3 2 4 4 5 5 6 5 7 4 8
8 6 10 8 4 3 9 2 4 1 2 2 3 1 4 3 5 4 6 5 7 7 8
5 10 10 4 10 4 1 2 2 3 1 4 4 5
10 2 1 6 6 8 5 7 7 5 6 1 2 1 3 2 4 2 5 3 6 3 7 5 8 8 9 4 10
2 8 10 1 2
6 10 7 8 6 2 8 1 2 2 3 1 4 1 5 5 6
2 8 7 1 2
1 9
9 6 6 3 3 8 9 2 7 10 1 2 1 3 2 4 2 5 5 6 6 7 5 8 5 9
10 9 5 7 6 1 10 6 7 5 3 1 2 2 3 1 4 1 5 5 6 2 7 1 8 1 9 1 10
8 2 9 7 8 6 8 4 4 1 2 1 3 1 4 2 5 1 6 2 7 7 8
5 7 4 6 8 6 1 2 2 3 3 4 3 5
4 9 1 9 4 1 2 2 3 3 4
1 7
1 2
2 2 10 1 2
9 2 9 10 10 8 6 3 8 4 1 2 1 3 2 4 2 5 5 6 5 7 3 8 3 9
10 1 1 5 2 7 9 4 4 7 9 1 2 1 3 2 4 2 5 2 6 6 7 5 8 4 9 5 10
7 5 7 1 8 4 4 5 1 2 2 3 2 4 1 5 3 6 3 7
6 7 7 6 9 7 5 1 2 2 3 3 4 2 5 5 6
8 9 3 9 1 3 2 4 9 1 2 2 3 3 4 1 5 2 6 5 7 3 8
1 9
1 8
7 2 6 6 3 5 9 1 1 2 1 3 2 4 3 5 1 6 6 7
9 10 6 2 10 10 1 4 5 10 1 2 1 3 1 4 4 5 1 6 5 7 3 8 6 9
4 10 7 7 10 1 2 1 3 1 4
9 10 4 9 4 5 3 5 3 10 1 2 2 3 2 4 2 5 2 6 1 7 2 8 8 9
1 2
2 9 1 1 2
6 8 7 5 9 4 7 1 2 2 3 3 4 4 5 4 6
9 9 10 4 1 2 7 9 7 2 1 2 2 3 3 4 2 5 1 6 6 7 1 8 2 9
4 4 8 4 9 1 2 2 3 3 4
3 1 1 7 1 2 2 3
8 10 2 4 8 3 1 1 8 1 2 1 3 1 4 4 5 4 6 4 7 5 8
4 6 4 2 7 1 2 1 3 1 4
5 10 10 5 9 2 1 2 1 3 3 4 3 5
10 6 5 5 7 7 6 8 6 4 2 1 2 2 3 1 4 2 5 3 6 5 7 7 8 1 9 2 10
2 3 1 1 2
1 7
1 3
4 1 2 5 3 1 2 1 3 1 4
10 4 1 6 7 2 1 6 9 3 8 1 2 2 3 2 4 1 5 5 6 5 7 4 8 5 9 8 10
10 10 10 1 5 5 3 1 4 5 2 1 2 1 3 3 4 3 5 3 6 1 7 6 8 4 9 6 10
7 6 2 4 5 2 8 10 1 2 1 3 2 4 3 5 2 6 1 7
7 5 2 3 1 1 1 1 1 2 1 3 2 4 1 5 2 6 3 7
5 5 4 4 4 6 1 2 1 3 1 4 1 5
4 8 2 9 7 1 2 2 3 2 4
`

type testCase struct {
	n       int
	weights []int64
	edges   [][2]int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesData, "\n")
	var cases []testCase
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			return nil, fmt.Errorf("invalid testcase line: %q", line)
		}
		idx := 0
		n64, err := strconv.ParseInt(fields[idx], 10, 64)
		if err != nil {
			return nil, err
		}
		n := int(n64)
		idx++
		if len(fields) < 1+n+n*2-2 { // rough sanity
			return nil, fmt.Errorf("not enough fields for n=%d in line: %q", n, line)
		}
		weights := make([]int64, n)
		for i := 0; i < n; i++ {
			w, err := strconv.ParseInt(fields[idx+i], 10, 64)
			if err != nil {
				return nil, err
			}
			weights[i] = w
		}
		idx += n
		edgeCount := n - 1
		edges := make([][2]int, edgeCount)
		for i := 0; i < edgeCount; i++ {
			u, err := strconv.Atoi(fields[idx+2*i])
			if err != nil {
				return nil, err
			}
			v, err := strconv.Atoi(fields[idx+2*i+1])
			if err != nil {
				return nil, err
			}
			edges[i] = [2]int{u, v}
		}
		cases = append(cases, testCase{n: n, weights: weights, edges: edges})
	}
	return cases, nil
}

// solve mirrors 1092F.go to compute the maximum value.
func solve(tc testCase) int64 {
	n := tc.n
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		a[i] = tc.weights[i-1]
	}
	graph := make([][]int, n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}
	parent := make([]int, n+1)
	depth := make([]int64, n+1)
	order := make([]int, 0, n)
	stack := make([]int, 0, n)
	stack = append(stack, 1)
	parent[1] = 0
	for i := 0; i < len(stack); i++ {
		x := stack[i]
		order = append(order, x)
		for _, y := range graph[x] {
			if y == parent[x] {
				continue
			}
			parent[y] = x
			depth[y] = depth[x] + 1
			stack = append(stack, y)
		}
	}
	val := make([]int64, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		x := order[i]
		val[x] = a[x]
		for _, y := range graph[x] {
			if y == parent[x] {
				continue
			}
			val[x] += val[y]
		}
	}
	dp := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dp[1] += depth[i] * a[i]
	}
	for _, x := range order {
		for _, y := range graph[x] {
			if y == parent[x] {
				continue
			}
			dp[y] = dp[x] + val[1] - 2*val[y]
		}
	}
	res := dp[1]
	for i := 2; i <= n; i++ {
		if dp[i] > res {
			res = dp[i]
		}
	}
	return res
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte(' ')
	for i, w := range tc.weights {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(w, 10))
	}
	for _, e := range tc.edges {
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(e[0]))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(e[1]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := buildInput(tc)
		want := strconv.FormatInt(solve(tc), 10)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
