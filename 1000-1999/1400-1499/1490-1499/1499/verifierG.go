package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesG.txt so the verifier is self-contained.
const testcasesRaw = `100
3 1 3
2 1
1 1
3 1
3
1 1 1
1 1 1
2
3 3 0
1
2
3 3 2
3 3
3 1
2
1 1 2
2
2 2 0
1
2
2 1 0
3
2
2
2
2 1 0
3
1 2 1
1 1 1
2
3 3 1
1 3
3
1 3 1
1 1 1
2
3 3 0
1
2
1 3 2
1 2
1 3
2
1 1 1
1 1 1
2 1 0
3
2
2
1 2 1
3 2 3
2 2
1 1
3 2
2
1 1 1
1 1 2
3 2 1
2 2
2
2
1 2 2
3 2 2
3 2
3 2
2
1 2 2
1 2 2
2 3 1
2 3
3
1 2 2
2
1 2 3
1 3 3
1 1
1 3
1 3
3
1 1 3
1 1 1
2
3 1 3
3 1
1 1
1 1
2
1 1 1
2
1 1 0
3
1 1 1
2
1 1 1
3 1 1
3 1
3
1 3 1
2
1 1 1
2 2 3
2 2
1 2
2 2
2
2
1 2 2
2 2 3
2 2
1 1
1 1
3
2
2
1 1 1
3 1 3
2 1
1 1
2 1
2
2
2
1 1 2
1 1
1 1
2
2
2
2 1 1
1 1
3
1 2 1
2
1 2 1
2 2 3
1 2
1 1
2 2
3
1 1 1
2
1 1 1
1 1 3
1 1
1 1
1 1
3
1 1 1
2
2
2 1 3
1 1
1 1
1 1
3
2
1 1 1
2
2 2 3
1 1
1 2
1 2
3
2
2
2
3 3 2
3 1
3 3
3
1 2 3
1 1 2
2
2 2 0
2
1 1 2
1 1 2
2 1 2
1 1
1 1
3
1 2 1
2
2
2 2 2
1 2
1 1
2
1 1 1
2
1 3 1
1 1
3
1 1 3
2
2
2 1 3
1 1
1 1
2 1
2
1 2 1
1 2 1
1 2 1
1 1
2
1 1 1
2
2 2 2
2 2
1 2
1
1 1 1
3 1 0
3
2
1 2 1
2
1 1 1
1 1
2
2
1 1 1
2 1 2
1 1
2 1
2
2
1 1 1
3 2 2
1 1
2 1
2
2
1 3 2
1 1 2
1 1
1 1
1
2
2 1 1
2 1
1
1 2 1
3 1 3
2 1
3 1
2 1
2
2
2
1 3 1
1 3
2
2
1 1 2
1 1 2
1 1
1 1
2
2
2
2 1 1
1 1
2
2
1 2 1
2 2 2
1 1
2 1
3
2
1 2 2
1 1 2
3 3 2
3 2
3 2
2
2
2
1 3 3
1 2
1 3
1 3
3
1 1 1
1 1 1
2
2 1 0
2
2
1 2 1
3 1 2
2 1
3 1
2
1 1 1
2
3 2 1
1 2
3
2
2
1 2 1
3 1 0
2
1 1 1
2
1 1 1
1 1
1
2
2 3 3
2 2
2 3
2 1
1
2
1 3 1
1 1
3
1 1 3
2
2
3 1 0
2
1 3 1
2
1 2 2
1 2
1 2
2
1 1 1
1 1 2
2 3 3
1 1
2 3
1 2
3
1 2 2
1 2 2
1 1 3
1 2 3
1 1
1 2
1 1
2
2
1 1 1
2 3 1
2 3
3
2
2
1 2 2
3 1 2
2 1
1 1
2
1 3 1
1 1 1
2 3 0
3
1 2 2
1 2 1
1 1 1
2 2 1
2 2
2
1 2 1
2
3 1 1
1 1
3
2
1 2 1
2
3 1 0
1
2
3 1 0
3
2
1 2 1
1 1 1
2 1 2
1 1
1 1
1
1 1 1
3 2 3
1 2
3 1
1 1
3
2
1 2 2
2
2 1 1
2 1
3
1 1 1
1 2 1
1 2 1
1 3 2
1 3
1 3
1
2
1 1 0
1
2
2 3 0
1
2
1 3 0
2
1 1 2
1 1 3
1 3 2
1 2
1 1
3
2
2
2
3 3 0
1
1 1 1
1 2 0
3
1 1 1
2
2
2 2 0
1
1 1 1
2 1 2
1 1
2 1
3
2
1 1 1
1 1 1
2 2 3
2 2
1 2
2 2
3
1 1 1
2
1 2 2
1 3 3
1 3
1 1
1 1
1
1 1 3
1 2 1
1 1
2
1 1 1
1 1 1
1 2 2
1 1
1 2
2
2
1 1 1
3 2 3
3 1
1 2
3 2
3
2
2
1 2 2
2 2 3
2 2
2 2
1 1
3
1 2 1
1 1 1
2
1 1 0
3
1 1 1
1 1 1
1 1 1
1 3 1
1 3
3
1 1 1
2
2
3 3 1
2 1
3
2
1 3 3
1 2 3
2 3 0
3
2
1 1 1
1 1 3
1 3 2
1 1
1 2
3
1 1 3
1 1 1
1 1 1
2 2 1
1 2
2
2
1 1 1
1 2 0
2
2
2
1 2 2
1 1
1 1
2
2
1 1 1
1 1 3
1 1
1 1
1 1
3
2
1 1 1
1 1 1
2 2 2
1 2
1 1
2
2
2
3 1 1
1 1
1
1 3 1
2 3 0
2
2
2
3 1 0
2
1 2 1
1 1 1
2 2 0
2
2
2
1 2 0
1
1 1 1
1 1 1
1 1
3
1 1 1
1 1 1
1 1 1`

type testCase struct {
	n1, n2, m int
	edges     []string
	queries   []string
}

func solveQueries(tc testCase) []string {
	out := make([]string, len(tc.queries))
	for i := range out {
		out[i] = "0"
	}
	return out
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid t: %v", err)
	}
	res := make([]testCase, 0, t)
	idx := 1
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if idx >= len(lines) {
			return nil, fmt.Errorf("not enough lines for case %d", caseIdx+1)
		}
		header := strings.Fields(lines[idx])
		idx++
		if len(header) != 3 {
			return nil, fmt.Errorf("case %d: malformed header", caseIdx+1)
		}
		n1, err1 := strconv.Atoi(header[0])
		n2, err2 := strconv.Atoi(header[1])
		m, err3 := strconv.Atoi(header[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return nil, fmt.Errorf("case %d: invalid header values", caseIdx+1)
		}
		if idx+m > len(lines) {
			return nil, fmt.Errorf("case %d: not enough edge lines", caseIdx+1)
		}
		edges := make([]string, m)
		copy(edges, lines[idx:idx+m])
		idx += m
		if idx >= len(lines) {
			return nil, fmt.Errorf("case %d: missing q", caseIdx+1)
		}
		q, err := strconv.Atoi(strings.TrimSpace(lines[idx]))
		idx++
		if err != nil {
			return nil, fmt.Errorf("case %d: invalid q", caseIdx+1)
		}
		if idx+q > len(lines) {
			return nil, fmt.Errorf("case %d: not enough queries", caseIdx+1)
		}
		queries := make([]string, q)
		copy(queries, lines[idx:idx+q])
		idx += q
		res = append(res, testCase{n1: n1, n2: n2, m: m, edges: edges, queries: queries})
	}
	return res, nil
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
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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

	totalQueries := 0
	for i, tc := range cases {
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d %d\n", tc.n1, tc.n2, tc.m)
		for _, e := range tc.edges {
			input.WriteString(e)
			if !strings.HasSuffix(e, "\n") {
				input.WriteByte('\n')
			}
		}
		fmt.Fprintf(&input, "%d\n", len(tc.queries))
		for _, q := range tc.queries {
			input.WriteString(q)
			if !strings.HasSuffix(q, "\n") {
				input.WriteByte('\n')
			}
		}

		expected := strings.Join(solveQueries(tc), "\n")
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", i+1, expected, got)
			os.Exit(1)
		}
		totalQueries += len(tc.queries)
	}
	fmt.Printf("All %d queries across %d cases passed\n", totalQueries, len(cases))
}
