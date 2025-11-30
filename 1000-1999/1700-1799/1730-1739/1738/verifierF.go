package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded test cases from testcasesF.txt. Each line: n m followed by m edge pairs.
const testcasesRaw = `5 4 1 4 1 2 3 4 3 5
2 1 1 2
2 0
2 1 1 2
4 1 2 3
1 0
3 0
5 10 2 4 1 2 3 4 1 5 1 4 2 3 4 5 2 5 1 3 3 5
1 0
4 2 2 4 1 2
4 2 1 2 3 4
2 1 1 2
4 5 2 4 1 2 3 4 2 3 1 3
5 0
4 4 2 3 2 4 1 3 3 4
5 9 2 4 1 2 3 4 1 5 2 3 4 5 2 5 1 3 3 5
4 2 1 3 1 4
5 9 2 4 3 4 1 5 1 4 2 3 4 5 2 5 1 3 3 5
5 4 2 3 1 3 3 4 1 4
1 0
2 1 1 2
4 4 2 3 2 4 1 3 1 4
4 2 2 3 3 4
1 0
1 0
4 4 2 3 1 3 3 4 1 4
4 2 2 3 1 4
1 0
3 0
1 0
3 3 2 3 1 2 1 3
2 1 1 2
2 0
4 6 2 4 1 2 3 4 1 4 2 3 1 3
1 0
5 6 2 4 1 2 3 4 2 3 2 5 3 5
3 3 2 3 1 2 1 3
5 3 2 4 1 2 1 3
4 3 1 2 3 4 1 4
3 3 2 3 1 2 1 3
1 0
1 0
3 2 2 3 1 3
3 3 2 3 1 2 1 3
3 2 2 3 1 3
2 0
2 0
1 0
3 0
4 0
4 6 2 4 1 2 3 4 1 4 2 3 1 3
4 0
3 3 2 3 1 2 1 3
2 0
4 1 3 4
1 0
5 2 2 3 1 3
2 1 1 2
1 0
1 0
1 0
3 3 2 3 1 2 1 3
5 2 1 3 1 4
3 2 2 3 1 2
4 1 1 4
2 1 1 2
1 0
4 5 2 4 1 2 3 4 2 3 1 3
2 1 1 2
1 0
4 0
3 0
5 4 4 5 2 5 3 4 1 5
5 0
5 3 1 2 1 3 3 4
1 0
3 2 2 3 1 3
1 0
2 1 1 2
4 1 2 3
4 2 1 2 1 3
2 1 1 2
2 0
1 0
4 1 2 3
4 2 2 4 1 3
1 0
1 0
5 8 2 4 1 5 1 4 2 3 4 5 2 5 1 3 3 5
1 0
3 1 1 2
1 0
2 0
3 2 2 3 1 2
4 3 2 3 3 4 1 4
1 0
3 0
1 0
3 1 1 2
4 5 2 4 3 4 1 4 2 3 1 3`

type testCase struct {
	n     int
	edges [][2]int
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 2 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d parse n: %v", idx+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d parse m: %v", idx+1, err)
		}
		if len(fields) != 2+2*m {
			return nil, fmt.Errorf("line %d expected %d numbers got %d", idx+1, 2+2*m, len(fields))
		}
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			u, err := strconv.Atoi(fields[2+2*i])
			if err != nil {
				return nil, fmt.Errorf("line %d parse u%d: %v", idx+1, i, err)
			}
			v, err := strconv.Atoi(fields[3+2*i])
			if err != nil {
				return nil, fmt.Errorf("line %d parse v%d: %v", idx+1, i, err)
			}
			edges[i] = [2]int{u, v}
		}
		cases = append(cases, testCase{n: n, edges: edges})
	}
	return cases, nil
}

// Embedded solver logic from 1738F.go: compute component colors.
func solve(tc testCase) []int {
	g := make([][]int, tc.n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	color := make([]int, tc.n+1)
	cur := 0
	q := make([]int, 0, tc.n)
	for i := 1; i <= tc.n; i++ {
		if color[i] != 0 {
			continue
		}
		cur++
		q = q[:0]
		q = append(q, i)
		color[i] = cur
		for len(q) > 0 {
			u := q[0]
			q = q[1:]
			for _, v := range g[u] {
				if color[v] == 0 {
					color[v] = cur
					q = append(q, v)
				}
			}
		}
	}
	return color[1:]
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
		want := solve(tc)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
		for _, e := range tc.edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}

		gotStr, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotParts := strings.Fields(strings.TrimSpace(gotStr))
		if len(gotParts) != tc.n {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\n", idx+1, tc.n, len(gotParts))
			os.Exit(1)
		}
		for i, part := range gotParts {
			val, err := strconv.Atoi(part)
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: parse output at pos %d: %v\n", idx+1, i, err)
				os.Exit(1)
			}
			if val != want[i] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %v\n", idx+1, want, gotParts)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
