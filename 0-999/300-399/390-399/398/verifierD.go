package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n       int
	m       int
	q       int
	online  []int
	edges   [][2]int
	queries []query
}

type query struct {
	op byte
	u  int
	v  int
}

// Embedded testcases (from testcasesD.txt) so the verifier is self contained.
const rawTestcases = `3 2 1|3|2 1 3|1 3|2 3|D 3 2
2 0 5|2|2 1|F 1|A 1 2|A 2 1|F 2|A 2 1
4 5 4|4|2 1 3 4|2 3|1 3|2 4|3 4|1 2|C 3|C 2|D 4 3|A 4 2
3 1 3|2|1 3|2 3|A 3 2|A 1 2|O 2
3 3 3|1|2|1 2|2 3|1 3|C 1|A 1 2|F 2
3 0 3|1|2|O 2|O 2|A 1 3
4 2 2|3|1 2 4|1 2|2 4|C 2|A 1 3
2 1 2|1|2|1 2|A 1 2|F 2
3 0 1|0||C 3
3 0 5|3|3 1 2|C 2|D 3 1|D 1 2|D 2 3|O 1
5 3 2|5|4 2 5 1 3|4 5|1 4|1 5|F 3|O 3
2 1 1|1|2|1 2|O 2
3 1 3|2|3 2|1 2|D 3 1|A 2 3|D 1 3
3 0 4|1|2|F 3|D 3 2|D 2 1|O 1
3 0 5|2|1 2|C 1|F 1|C 2|C 2|C 2
2 0 1|2|1 2|C 2
3 0 3|0||A 2 1|O 3|D 2 1
4 3 2|3|4 1 3|3 4|1 3|2 3|D 4 1|F 3
5 3 1|4|2 3 1 4|1 3|1 5|2 3|O 3
2 0 4|2|2 1|C 1|D 2 1|F 1|O 1
4 4 4|4|3 4 2 1|1 2|2 4|3 4|1 4|A 2 3|O 2|C 2|A 3 1
4 5 5|0||3 4|1 4|2 4|1 3|2 3|A 2 3|O 1|C 1|A 1 2|F 4
2 0 2|2|2 1|C 2|D 1 2
3 2 3|1|3|1 2|2 3|D 3 2|A 2 3|D 3 1
4 1 5|3|3 2 4|2 3|D 3 4|A 1 3|D 2 4|D 1 2|F 1
4 1 1|1|4|1 3|F 1
3 0 3|0||O 3|A 3 1|F 3
5 1 1|4|3 1 5 4|2 4|F 3
5 4 3|3|3 2 4|1 3|1 5|2 3|3 5|F 1|C 1|F 2
2 1 5|0||1 2|D 2 1|F 1|O 2|C 2|A 2 1
3 2 5|2|2 3|1 3|2 3|A 1 2|C 2|C 2|A 1 2|O 2
3 0 2|0||A 1 2|O 1
4 4 1|0||1 4|2 4|3 4|1 3|F 3
4 3 5|3|1 4 2|2 4|1 3|3 4|F 4|A 2 1|O 4|D 3 2|D 4 2
3 2 3|0||2 3|1 3|A 3 2|F 3|A 1 2
2 0 3|2|1 2|O 2|O 2|C 1
3 0 2|0||A 3 1|O 3
3 0 2|1|2|A 3 2|F 1
3 0 1|1|3|D 1 2
3 3 4|3|3 2 1|1 2|2 3|1 3|O 1|F 1|A 1 3|F 2
3 2 4|2|3 1|1 3|2 3|A 1 2|O 1|O 1|A 3 2
5 0 3|5|4 2 1 5 3|F 1|A 3 1|O 4
3 0 2|1|3|O 2|D 1 2
2 0 2|2|2 1|D 2 1|F 2
4 2 5|2|1 3|1 2|2 4|O 1|C 2|F 2|C 2|C 1
4 0 1|4|2 3 1 4|F 4
4 4 1|1|3|3 4|1 4|1 3|2 4|O 2
5 5 1|3|1 3 5|1 3|3 5|3 4|2 3|1 2|A 5 4
2 0 5|2|2 1|F 1|O 1|F 2|F 1|C 2
4 5 3|1|1|1 3|2 3|1 4|1 2|3 4|F 1|A 4 2|C 3
3 2 5|3|1 2 3|2 3|1 2|C 3|C 1|A 3 2|A 1 3|F 2
4 2 2|4|4 3 2 1|2 4|2 3|A 2 1|D 1 3
2 0 1|0||C 1
3 1 3|3|2 1 3|1 2|A 3 1|A 1 2|F 1
5 0 4|1|4|O 3|A 1 5|C 2|D 3 5
2 1 2|0||1 2|D 2 1|O 1
4 0 2|1|4|A 4 3|O 3
5 4 1|5|4 5 3 2 1|1 3|1 5|2 3|2 4|C 2
5 1 5|2|4 5|2 4|F 2|A 3 2|A 5 4|C 2|D 4 5
4 4 5|0||1 3|1 2|2 4|1 4|C 2|O 4|C 3|F 1|O 2
4 1 3|2|4 3|3 4|O 4|O 2|F 4
2 0 2|0||O 1|C 1
3 0 5|3|2 3 1|C 1|A 3 1|C 2|A 1 2|A 3 1
5 0 2|5|2 1 5 3 4|F 2|O 2
5 4 5|3|4 3 5|2 4|3 5|1 3|4 5|A 1 3|C 2|F 4|C 2|F 5
3 0 5|0||O 3|D 3 1|C 1|A 2 1|F 1
4 4 4|2|3 1|1 3|3 4|2 4|1 4|D 3 4|F 4|F 2|F 3
2 0 2|2|2 1|F 1|A 2 1
4 2 1|3|4 2 3|2 4|3 4|C 4
5 0 3|2|5 1|D 2 1|A 2 3|D 5 2
3 3 4|2|2 3|1 2|2 3|1 3|F 3|O 3|C 1|C 1
2 1 3|1|1|1 2|F 2|O 2|D 2 1
4 3 3|3|1 2 4|1 4|1 2|3 4|D 3 2|C 1|F 2
4 4 5|1|1|3 4|1 3|1 4|1 2|C 2|O 4|D 3 2|F 1|O 2
5 2 5|3|1 3 2|2 5|1 5|A 1 3|D 3 1|F 5|O 1|D 5 3
5 5 1|3|4 5 3|1 5|1 2|2 4|2 5|3 5|D 5 2
4 5 5|3|3 4 2|2 3|3 4|1 3|1 2|1 4|A 3 1|F 2|C 3|C 1|A 4 2
4 5 4|4|2 3 4 1|1 3|1 2|3 4|2 4|1 4|A 4 1|F 4|O 4|C 2
2 1 5|0||1 2|F 2|D 2 1|F 1|C 1|D 2 1
4 5 3|0||2 3|1 3|3 4|1 2|2 4|D 4 1|F 3|D 2 4
4 3 5|1|2|1 2|2 4|1 4|F 1|A 2 4|F 3|F 2|A 4 3
2 0 3|2|2 1|D 1 2|F 1|C 1
4 0 1|1|2|F 1
4 2 3|3|4 1 3|1 2|1 3|A 1 4|O 4|A 1 4
4 0 3|4|1 3 2 4|C 4|A 1 4|F 2
5 4 3|4|3 4 5 1|3 4|2 4|1 2|3 5|D 3 2|C 5|D 1 3
2 0 4|0||D 1 2|F 1|F 1|O 2
3 2 2|2|1 3|1 2|2 3|D 2 3|F 1
5 3 4|5|5 4 3 2 1|1 2|1 4|4 5|D 2 4|D 2 4|F 4|O 5
3 1 5|0||1 2|D 3 2|C 3|O 3|C 2|F 1
2 0 4|0||F 1|O 2|A 1 2|F 1
4 5 5|3|3 1 2|1 4|1 2|2 4|3 4|1 3|D 2 3|O 1|F 1|A 4 2|A 2 1
4 5 1|2|1 2|1 3|1 2|2 4|1 4|3 4|D 3 4
5 3 3|0||2 3|2 4|1 4|A 2 5|F 3|A 4 2
3 3 1|1|2|1 3|2 3|1 2|A 1 2
4 4 4|3|2 1 3|1 4|2 3|1 3|2 4|F 2|A 4 1|F 3|F 4
2 1 5|1|2|1 2|F 2|F 2|D 1 2|A 2 1|C 1
2 0 1|2|2 1|C 1
5 4 2|3|4 5 3|4 5|1 2|1 4|2 5|F 1|D 2 5
4 2 1|0||3 4|1 3|D 4 3`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(rawTestcases, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 3 {
			return nil, fmt.Errorf("line %d: expected at least 3 parts", idx+1)
		}
		header := strings.Fields(parts[0])
		if len(header) != 3 {
			return nil, fmt.Errorf("line %d: expected 3 numbers in header got %d", idx+1, len(header))
		}
		n, err := strconv.Atoi(header[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		m, err := strconv.Atoi(header[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %w", idx+1, err)
		}
		q, err := strconv.Atoi(header[2])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse q: %w", idx+1, err)
		}

		oStr := strings.TrimSpace(parts[1])
		o := 0
		if oStr != "" {
			o, err = strconv.Atoi(oStr)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse o: %w", idx+1, err)
			}
		}

		onlineFields := strings.Fields(parts[2])
		if len(onlineFields) != o {
			return nil, fmt.Errorf("line %d: expected %d online nodes got %d", idx+1, o, len(onlineFields))
		}
		online := make([]int, o)
		for i, f := range onlineFields {
			val, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse online node %d: %w", idx+1, i, err)
			}
			online[i] = val
		}

		expectedParts := 3 + m + q
		if len(parts) != expectedParts {
			return nil, fmt.Errorf("line %d: expected %d parts got %d", idx+1, expectedParts, len(parts))
		}

		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			f := strings.Fields(parts[3+i])
			if len(f) != 2 {
				return nil, fmt.Errorf("line %d: edge %d malformed", idx+1, i+1)
			}
			u, err := strconv.Atoi(f[0])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse edge %d u: %w", idx+1, i+1, err)
			}
			v, err := strconv.Atoi(f[1])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse edge %d v: %w", idx+1, i+1, err)
			}
			edges[i] = [2]int{u, v}
		}

		queries := make([]query, q)
		for i := 0; i < q; i++ {
			f := strings.Fields(parts[3+m+i])
			if len(f) < 2 {
				return nil, fmt.Errorf("line %d: query %d malformed", idx+1, i+1)
			}
			op := f[0][0]
			switch op {
			case 'O', 'F', 'C':
				if len(f) != 2 {
					return nil, fmt.Errorf("line %d: query %d expected 1 arg", idx+1, i+1)
				}
				u, err := strconv.Atoi(f[1])
				if err != nil {
					return nil, fmt.Errorf("line %d: parse query %d arg: %w", idx+1, i+1, err)
				}
				queries[i] = query{op: op, u: u}
			case 'A', 'D':
				if len(f) != 3 {
					return nil, fmt.Errorf("line %d: query %d expected 2 args", idx+1, i+1)
				}
				u, err := strconv.Atoi(f[1])
				if err != nil {
					return nil, fmt.Errorf("line %d: parse query %d u: %w", idx+1, i+1, err)
				}
				v, err := strconv.Atoi(f[2])
				if err != nil {
					return nil, fmt.Errorf("line %d: parse query %d v: %w", idx+1, i+1, err)
				}
				queries[i] = query{op: op, u: u, v: v}
			default:
				return nil, fmt.Errorf("line %d: query %d has unknown op %c", idx+1, i+1, op)
			}
		}

		cases = append(cases, testCase{n: n, m: m, q: q, online: online, edges: edges, queries: queries})
	}
	return cases, nil
}

func solveCase(tc testCase) string {
	status := make([]bool, tc.n+1)
	for _, v := range tc.online {
		if v >= 0 && v <= tc.n {
			status[v] = true
		}
	}

	adj := make([]map[int]struct{}, tc.n+1)
	for i := 0; i <= tc.n; i++ {
		adj[i] = make(map[int]struct{})
	}
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		adj[u][v] = struct{}{}
		adj[v][u] = struct{}{}
	}

	var outputs []string
	for _, q := range tc.queries {
		switch q.op {
		case 'O':
			status[q.u] = true
		case 'F':
			status[q.u] = false
		case 'A':
			adj[q.u][q.v] = struct{}{}
			adj[q.v][q.u] = struct{}{}
		case 'D':
			delete(adj[q.u], q.v)
			delete(adj[q.v], q.u)
		case 'C':
			cnt := 0
			for v := range adj[q.u] {
				if status[v] {
					cnt++
				}
			}
			outputs = append(outputs, strconv.Itoa(cnt))
		}
	}

	return strings.Join(outputs, "\n")
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.q))
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.online)))
	if len(tc.online) > 0 {
		for i, v := range tc.online {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for _, q := range tc.queries {
		switch q.op {
		case 'O', 'F', 'C':
			sb.WriteString(fmt.Sprintf("%c %d\n", q.op, q.u))
		case 'A', 'D':
			sb.WriteString(fmt.Sprintf("%c %d %d\n", q.op, q.u, q.v))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expected := solveCase(tc)
		input := buildInput(tc)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx+1, err, string(out))
			os.Exit(1)
		}

		got := strings.TrimSpace(string(out))
		if got != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected:\n%s\n got:\n%s\n", idx+1, strings.TrimSpace(expected), got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
