package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesRaw = `2 1 0 1 1 0
1 0 0
4 4 1 0 0 1 2 1 2 0 1 0 3 2
5 3 0 0 1 1 0 4 2 2 0 4 1
4 3 1 1 1 0 2 1 2 0 3 2
5 8 1 1 0 0 0 3 1 3 2 2 1 4 3 4 2 1 0 2 0 4 0
3 1 1 1 0 2 0
3 0 1 0 1
2 1 1 1 1 0
2 1 0 0 1 0
3 0 0 1 1
6 6 1 0 0 1 1 1 1 0 5 1 3 1 4 1 3 2 2 1
1 0 0
5 6 0 1 1 0 1 4 1 2 1 3 2 2 0 1 0 4 3
2 1 0 0 1 0
1 0 1
2 0 0 0
6 10 1 1 0 0 0 0 4 3 4 2 1 0 2 0 4 0 3 1 3 2 3 0 5 3 5 1
6 9 1 1 0 0 1 0 5 1 2 1 3 2 4 2 5 0 1 0 5 2 4 0 3 0
6 15 0 0 1 1 1 0 5 1 5 0 4 0 4 2 4 3 3 1 5 2 2 1 3 0 4 1 5 4 5 3 2 0 1 0 3 2
5 1 1 0 0 1 1 4 2
4 2 0 1 0 1 2 0 1 0
2 0 0 0
2 0 1 0
4 5 1 1 0 0 3 0 2 1 3 2 2 0 1 0
4 4 0 0 0 0 2 1 2 0 3 0 1 0
4 3 0 1 1 0 3 2 3 0 2 1
3 0 0 1 1
5 4 0 1 0 0 0 2 1 3 1 2 0 3 0
5 3 0 0 0 0 0 1 0 2 0 4 2
6 10 1 0 0 0 0 0 5 1 5 0 5 2 3 0 4 3 5 3 1 0 2 1 4 1 5 4
3 3 1 0 0 1 0 2 0 2 1
4 0 1 1 0 0
4 4 1 1 1 0 3 1 3 0 2 0 1 0
2 0 1 0
6 0 1 1 0 0 0 0
2 1 0 1 1 0
6 11 1 1 0 0 1 0 5 4 4 1 5 1 4 0 4 2 5 0 2 1 3 2 5 3 3 0 2 0
5 9 1 0 0 1 0 3 0 2 1 1 0 2 0 4 3 4 2 3 1 4 0 4 1
5 1 0 0 0 0 1 4 3
3 2 0 1 1 2 0 1 0
4 5 0 0 1 1 2 1 3 2 1 0 3 1 2 0
1 0 0
1 0 0
2 1 1 1 1 0
6 13 1 0 1 0 0 0 3 0 3 1 4 0 2 1 4 2 5 2 1 0 4 1 3 2 4 3 5 4 2 0 5 3
3 3 0 1 1 2 0 2 1 1 0
5 9 0 1 1 0 1 4 0 4 1 3 1 4 2 4 3 3 0 3 2 2 1 1 0
1 0 1
2 1 1 1 1 0
1 0 0
4 0 0 0 1 0
5 9 0 0 0 0 0 4 1 1 0 4 3 4 2 2 1 4 0 2 0 3 0 3 2
1 0 0
6 14 1 1 0 0 0 0 5 1 4 3 3 2 4 0 5 4 5 2 1 0 4 1 5 3 2 0 5 0 2 1 4 2 3 1
2 0 1 0
2 0 1 1
2 0 1 1
6 14 1 1 0 1 1 1 2 1 5 0 5 2 5 3 2 0 4 0 4 2 3 1 4 1 5 1 3 2 4 3 1 0 5 4
2 0 1 0
6 8 0 0 1 1 1 1 1 0 3 2 4 3 5 0 5 4 5 3 4 0 4 1
4 6 1 0 0 0 3 0 3 2 3 1 2 0 2 1 1 0
3 0 0 1 1
3 2 0 0 0 2 0 2 1
3 0 1 0 1
6 14 1 0 0 1 1 0 4 2 5 3 4 1 5 1 3 1 4 0 2 1 5 2 3 0 5 4 5 0 2 0 4 3 3 2
5 0 0 0 1 1 1
4 5 0 1 0 1 3 1 2 0 2 1 3 0 1 0
6 2 1 0 1 0 0 0 5 1 2 0
2 1 0 1 1 0
1 0 1
4 5 1 1 1 1 2 1 3 1 3 2 1 0 3 0
6 9 0 0 0 0 1 0 3 0 2 1 1 0 4 1 5 2 5 4 3 2 2 0 3 1
4 2 0 1 0 1 1 0 3 2
4 6 1 0 0 1 2 0 2 1 3 2 1 0 3 0 3 1
2 1 1 1 1 0
2 1 0 0 1 0
2 1 1 1 1 0
6 4 0 1 1 0 0 1 5 1 4 3 4 0 3 2
2 1 0 1 1 0
6 7 1 1 1 0 0 0 5 0 2 1 4 1 4 3 5 1 2 0 5 3
6 14 0 0 0 0 1 0 5 3 3 2 5 0 2 0 5 4 3 1 5 1 4 2 4 3 1 0 4 0 5 2 2 1 4 1
2 1 0 1 1 0
2 0 0 1
2 0 1 1
6 9 0 1 1 1 0 1 5 0 4 1 2 0 5 1 5 3 5 4 3 0 3 1 2 1
2 0 0 0
4 6 0 1 1 1 3 0 2 1 3 1 2 0 1 0 3 2
2 1 1 0 1 0
2 1 0 1 1 0
6 2 1 1 1 0 0 0 4 2 3 0
5 3 1 1 1 0 0 4 3 2 1 4 1
2 1 0 1 1 0
5 7 1 0 1 0 1 4 1 2 0 4 0 4 2 2 1 4 3 3 0
5 2 0 0 0 1 1 4 0 4 1
3 1 1 0 0 2 0
1 0 0
3 0 1 0 0
6 7 1 0 0 0 0 0 3 1 5 2 4 1 5 4 4 2 1 0 4 0
3 2 1 0 1 2 0 1 0`

type testCase struct {
	n    int
	m    int
	typ  []int
	edge [][2]int
}

func solveCase(tc testCase) (int, error) {
	n := tc.n
	m := tc.m
	if len(tc.typ) != n {
		return 0, fmt.Errorf("type length mismatch: got %d want %d", len(tc.typ), n)
	}
	if len(tc.edge) != m {
		return 0, fmt.Errorf("edge count mismatch: got %d want %d", len(tc.edge), m)
	}
	adj := make([][]int, n)
	indeg := make([]int, n)
	for _, e := range tc.edge {
		t1, t2 := e[0], e[1]
		if t1 < 0 || t1 >= n || t2 < 0 || t2 >= n {
			return 0, fmt.Errorf("edge out of range: %v", e)
		}
		adj[t2] = append(adj[t2], t1)
		indeg[t1]++
	}

	qCPU := make([]int, 0)
	qGPU := make([]int, 0)
	for i := 0; i < n; i++ {
		if indeg[i] == 0 {
			if tc.typ[i] == 0 {
				qCPU = append(qCPU, i)
			} else {
				qGPU = append(qGPU, i)
			}
		}
	}

	processed := 0
	calls := 0
	for processed < n {
		for len(qCPU) > 0 {
			u := qCPU[0]
			qCPU = qCPU[1:]
			processed++
			for _, v := range adj[u] {
				indeg[v]--
				if indeg[v] == 0 {
					if tc.typ[v] == 0 {
						qCPU = append(qCPU, v)
					} else {
						qGPU = append(qGPU, v)
					}
				}
			}
		}
		if processed >= n {
			break
		}
		if len(qGPU) > 0 {
			calls++
			for len(qGPU) > 0 {
				u := qGPU[0]
				qGPU = qGPU[1:]
				processed++
				for _, v := range adj[u] {
					indeg[v]--
					if indeg[v] == 0 {
						if tc.typ[v] == 0 {
							qCPU = append(qCPU, v)
						} else {
							qGPU = append(qGPU, v)
						}
					}
				}
			}
		}
	}
	return calls, nil
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		p := 0
		n, err := strconv.Atoi(fields[p])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		p++
		m, err := strconv.Atoi(fields[p])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %w", idx+1, err)
		}
		p++
		if len(fields) != 2+n+2*m {
			return nil, fmt.Errorf("line %d: expected %d fields got %d", idx+1, 2+n+2*m, len(fields))
		}
		typ := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[p])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse type %d: %w", idx+1, i+1, err)
			}
			typ[i] = v
			p++
		}
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			a, err := strconv.Atoi(fields[p])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse edge %d a: %w", idx+1, i+1, err)
			}
			b, err := strconv.Atoi(fields[p+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse edge %d b: %w", idx+1, i+1, err)
			}
			edges[i] = [2]int{a, b}
			p += 2
		}
		cases = append(cases, testCase{n: n, m: m, typ: typ, edge: edges})
	}
	return cases, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		expectVal, err := solveCase(tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d solve error: %v\n", idx+1, err)
			os.Exit(1)
		}
		expect := strconv.Itoa(expectVal)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for i, v := range tc.typ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, e := range tc.edge {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		input := sb.String()
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\nInput:\n%s\n", idx+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
