package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod int64 = 998244353

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesRaw = `
100
2 1 1 1 2 1
3 1 5 2 1 1 2 1
3 2 5 5 3 1 3 1 2 2 3 2 2 2
3 2 2 1 3 1
4 4 5 4 2 3 3 4 4 2 4 4
3 4 1 4 2 2 3 4 1 2 1 3
4 4 1 5 1 1 4 1 3 1 1 1 3 3
2 2 5 3 1 1 2 2 2 2
2 3 4 4 1 2 1 3 2 1 1 1
4 3 4 4 4 2 2 1 4 1 1 2
1 2 1 2 1 1 1 1
3 1 5 3 2 1 2 1 3 1
4 4 3 4 2 3 3 1 1 1 2 3
1 3 1 2 1 2 1 3
4 2 4 2 3 1 1 2
3 4 3 3 1 1 1 3 1 2
3 2 2 3 3 2 2 2 2 2
3 2 1 3 3 1 1 2 2 2
1 4 4 4 1 1 1 1 1 1 1 4
2 4 5 1 2 4
1 2 2 4 1 1 1 2 1 2 1 2
2 3 4 5 2 3 2 2 1 1 1 3 1 1
4 2 2 3 1 2 1 1 4 2
1 3 2 5 1 1 1 1 1 3 1 1 1 3
4 3 5 2 2 1 3 2
1 2 4 2 1 1 1 1
4 2 2 2 2 2 3 1
4 1 5 1 3 1
4 3 1 2 2 3 4 3
4 3 1 3 2 3 4 3 2 2
3 4 1 2 1 3 3 2
3 4 2 3 2 2 2 1 2 1
2 3 1 4 1 1 1 1 1 1 2 1
1 3 4 2 1 2 1 1
3 3 2 3 3 3 2 2 1 2
3 4 1 5 1 4 2 2 1 2 3 2 1 3
2 2 2 1 1 1
1 4 5 1 1 2
1 3 3 4 1 3 1 2 1 2 1 3
2 4 2 5 2 2 2 3 2 4 2 1 2 3
4 1 5 5 3 1 4 1 1 1 3 1 1 1
3 1 3 5 3 1 3 1 3 1 2 1 1 1
3 2 5 5 3 2 1 1 2 1 3 1 3 1
3 3 3 3 2 2 2 3 2 1
1 2 5 1 1 1
3 1 4 3 3 1 1 1 2 1
2 1 2 5 1 1 2 1 2 1 1 1 2 1
3 3 2 4 1 1 3 1 3 3 3 2
3 4 1 5 2 2 1 3 3 2 3 2 3 3
4 1 2 5 2 1 2 1 1 1 2 1 1 1
2 4 1 4 2 3 2 4 2 3 1 3
1 4 1 3 1 4 1 4 1 1
4 3 5 2 3 2 1 2
4 2 1 1 3 2
1 1 2 1 1 1
3 1 3 4 1 1 2 1 3 1 2 1
1 2 2 2 1 1 1 2
2 3 5 1 1 3
2 3 2 3 2 3 2 2 2 1
4 2 1 5 2 1 3 1 2 2 2 2 2 2
2 4 5 1 1 3
1 1 1 4 1 1 1 1 1 1 1 1
1 2 2 1 1 1
1 4 1 3 1 1 1 3 1 4
1 4 3 5 1 3 1 4 1 1 1 3 1 4
1 3 1 1 1 2
1 2 4 5 1 2 1 2 1 1 1 2 1 1
1 1 3 3 1 1 1 1 1 1
1 4 3 2 1 2 1 2
2 3 1 4 2 1 1 1 2 2 1 1
3 2 2 4 2 2 3 1 2 1 1 1
4 2 2 1 1 2
1 1 1 1 1 1
4 1 4 4 3 1 2 1 2 1 1 1
4 3 1 4 4 1 4 2 2 2 2 3
4 4 5 1 3 3
3 1 5 5 2 1 2 1 3 1 2 1 3 1
2 4 1 4 2 4 1 4 2 2 1 3
4 4 2 5 3 1 3 3 2 3 4 1 2 3
2 3 3 4 1 2 2 1 2 2 2 3
2 3 2 2 1 3 2 1
1 4 5 4 1 1 1 4 1 2 1 2
2 2 1 5 1 2 2 2 2 1 2 2 1 2
3 2 4 5 2 1 2 2 2 2 3 2 2 2
1 2 2 1 1 1
1 3 5 2 1 1 1 1
1 3 4 3 1 3 1 2 1 3
4 2 2 5 1 1 1 1 4 1 4 2 3 1
2 4 2 4 1 4 1 4 2 2 2 1
4 1 3 5 4 1 4 1 2 1 3 1 3 1
4 1 1 4 4 1 1 1 4 1 1 1
2 3 1 2 1 3 1 3
2 1 1 4 2 1 2 1 2 1 1 1
2 4 3 3 1 1 2 3 2 1
1 3 5 4 1 3 1 2 1 1 1 1
4 3 3 4 2 1 4 2 4 1 3 3
3 4 3 3 2 3 1 4 2 1
4 2 3 2 2 2 4 1
1 4 4 4 1 2 1 3 1 2 1 1
1 4 3 4 1 4 1 1 1 1 1 2
`

type testCase struct {
	n int
	m int
	k int
	q int
	x []int
	y []int
}

func solveCase(tc testCase) int64 {
	rowUsed := make([]bool, tc.n+1)
	colUsed := make([]bool, tc.m+1)
	cntRow, cntCol := 0, 0
	ans := int64(1)
	for i := tc.q - 1; i >= 0; i-- {
		if cntRow == tc.n || cntCol == tc.m {
			break
		}
		x, y := tc.x[i], tc.y[i]
		if !rowUsed[x] || !colUsed[y] {
			ans = ans * int64(tc.k) % mod
			if !rowUsed[x] {
				rowUsed[x] = true
				cntRow++
			}
			if !colUsed[y] {
				colUsed[y] = true
				cntCol++
			}
		}
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	cases := make([]testCase, 0, t)
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 4 {
			return nil, fmt.Errorf("line %d: not enough fields", i+1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", i+1, err)
		}
		m, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %v", i+1, err)
		}
		k, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse k: %v", i+1, err)
		}
		q, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse q: %v", i+1, err)
		}
		expected := 4 + 2*q
		if len(parts) != expected {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", i+1, expected, len(parts))
		}
		x := make([]int, q)
		y := make([]int, q)
		for j := 0; j < q; j++ {
			vx, err := strconv.Atoi(parts[4+2*j])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse x[%d]: %v", i+1, j, err)
			}
			vy, err := strconv.Atoi(parts[5+2*j])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse y[%d]: %v", i+1, j, err)
			}
			x[j] = vx
			y[j] = vy
		}
		cases = append(cases, testCase{n: n, m: m, k: k, q: q, x: x, y: y})
	}
	if len(cases) != t {
		return nil, fmt.Errorf("expected %d cases got %d", t, len(cases))
	}
	return cases, nil
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
		expected := solveCase(tc) % mod

		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.n, tc.m, tc.k, tc.q))
		for j := 0; j < tc.q; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", tc.x[j], tc.y[j]))
		}

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotVal, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if err != nil || gotVal%mod != expected {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
