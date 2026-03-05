package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded test cases: all values satisfy 1 ≤ k ≤ 100000.
const testcasesAData = `
1
2
3
4
5
6
7
8
9
10
15
20
25
50
99
100
200
500
999
1000
1234
2000
3000
4000
4999
5000
7777
9999
10000
12345
15000
19999
20000
23456
29999
30000
34567
39999
40000
45678
49999
50000
54321
59999
60000
65432
69999
70000
75000
79999
80000
85000
89999
90000
95000
99999
100000
161
330
561
715
1001
1540
3003
4465
5985
7315
8001
11175
12341
14950
17550
19250
22100
24310
27405
32509
35001
38760
43758
48620
53130
58905
63756
66045
71253
74613
77520
83521
86526
88560
91390
94143
97251
98010
`

// generateGraph implements the logic from 232A.go, returning the constructed graph.
func generateGraph(k int) (int, [][]bool) {
	const maxN = 101
	adj := make([][]bool, maxN)
	for i := range adj {
		adj[i] = make([]bool, maxN)
	}
	n := 1
	for k > 0 {
		m := n
		for m*(m-1)/2 > k {
			m--
		}
		for i := 0; i < m; i++ {
			adj[i][n] = true
			adj[n][i] = true
		}
		k -= m * (m - 1) / 2
		n++
	}
	graph := make([][]bool, n)
	for i := 0; i < n; i++ {
		graph[i] = make([]bool, n)
		copy(graph[i], adj[i][:n])
	}
	return n, graph
}

func countTriangles(adj [][]bool) int {
	n := len(adj)
	cnt := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if !adj[i][j] {
				continue
			}
			for k := j + 1; k < n; k++ {
				if adj[i][k] && adj[j][k] {
					cnt++
				}
			}
		}
	}
	return cnt
}

func runCase(bin string, k int) error {
	input := fmt.Sprintf("%d\n", k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(out.String())))
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	n64, err := strconv.ParseInt(strings.TrimSpace(scanner.Text()), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse n: %v", err)
	}
	n := int(n64)
	if n < 3 || n > 100 {
		return fmt.Errorf("invalid n %d", n)
	}
	adj := make([][]bool, n)
	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("not enough lines for adjacency matrix")
		}
		line := strings.TrimSpace(scanner.Text())
		if len(line) != n {
			return fmt.Errorf("line %d length mismatch", i+1)
		}
		adj[i] = make([]bool, n)
		for j := 0; j < n; j++ {
			c := line[j]
			if c != '0' && c != '1' {
				return fmt.Errorf("invalid char at row %d col %d", i+1, j+1)
			}
			if i == j && c != '0' {
				return fmt.Errorf("self loop at row %d", i+1)
			}
			adj[i][j] = c == '1'
		}
	}
	if scanner.Scan() {
		extra := strings.TrimSpace(scanner.Text())
		if extra != "" {
			return fmt.Errorf("extra output: %q", extra)
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if adj[i][j] != adj[j][i] {
				return fmt.Errorf("asymmetry at %d,%d", i, j)
			}
		}
	}
	if countTriangles(adj) != k {
		return fmt.Errorf("expected %d triangles", k)
	}
	nRef, _ := generateGraph(k)
	if n > nRef {
		return fmt.Errorf("n (%d) larger than reference construction (%d)", n, nRef)
	}
	return nil
}

func parseTestcases() ([]int, error) {
	sc := bufio.NewScanner(strings.NewReader(testcasesAData))
	var cases []int
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		val, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("parse testcase: %w", err)
		}
		cases = append(cases, val)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, k := range testcases {
		if err := runCase(bin, k); err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
