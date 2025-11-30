package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `
4 3 8 7 2 6 4 5
3 1 6 4
3 1 4 2
4 4 1 6 5 8 7 4 3 2
4 4 4 2 3 6 7 8 1 5
1 1 1 2
2 2 2 1 4 3
2 2 4 3 1 2
4 0
1 1 2 1
3 3 3 6 2 4 5 1
3 1 4 5
4 0
3 1 1 6
4 3 4 8 2 7 5 3
2 0
1 0
2 1 2 3
2 2 4 3 1 2
1 0
1 1 2 1
3 0
1 1 2 1
1 0
1 1 2 1
2 0
2 0
4 1 8 1
1 0
3 2 2 6 3 4
1 1 1 2
3 2 6 1 5 2
4 0
4 2 1 4 7 2
4 3 6 4 7 3 8 5
2 2 1 3 2 4
1 1 1 2
1 0
1 1 2 1
4 4 4 3 2 6 1 7 5 8
3 3 5 6 3 2 4 1
2 0
3 0
2 0
2 0
1 1 1 2
1 1 2 1
3 0
2 2 1 2 4 3
3 0
4 2 7 6 1 8
2 1 1 4
1 1 2 1
1 1 2 1
2 2 1 4 2 3
2 1 3 4
1 1 2 1
1 1 1 2
3 2 6 2 3 5
4 2 3 5 2 8
3 2 1 4 3 5
1 0
2 2 3 4 1 2
1 0
3 3 3 1 6 2 5 4
3 1 6 1
3 0
1 1 1 2
3 1 4 5
1 0
4 2 7 6 8 2
4 4 2 1 7 5 4 3 6 8
2 0
3 0
4 4 2 3 8 1 4 7 6 5
1 0
1 1 2 1
4 2 1 6 5 8
1 0
4 1 3 1
1 0
3 2 5 6 2 4
4 3 8 3 4 5 6 2
1 0
4 1 6 3
1 1 1 2
4 0
3 0
4 0
3 1 2 6
1 1 1 2
3 1 2 5
4 2 4 3 8 2
4 4 6 3 7 8 4 2 1 5
4 0
1 1 1 2
3 0
3 2 3 5 4 6
2 0
1 1 1 2

`

type testCase struct {
	n     int
	k     int
	pairs [][2]int
}

func solveCase(tc testCase) int {
	pairs := make([][2]int, 0, tc.n)
	used := make([]bool, 2*tc.n+1)
	for _, p := range tc.pairs {
		x, y := p[0], p[1]
		if x > y {
			x, y = y, x
		}
		pairs = append(pairs, [2]int{x, y})
		used[x] = true
		used[y] = true
	}
	unused := make([]int, 0, 2*(tc.n-tc.k))
	for i := 1; i <= 2*tc.n; i++ {
		if !used[i] {
			unused = append(unused, i)
		}
	}
	m := tc.n - tc.k
	for i := 0; i < m; i++ {
		x := unused[i]
		y := unused[i+m]
		if x > y {
			x, y = y, x
		}
		pairs = append(pairs, [2]int{x, y})
	}
	// sort by first element
	for i := 1; i < len(pairs); i++ {
		j := i
		for j > 0 && pairs[j][0] < pairs[j-1][0] {
			pairs[j], pairs[j-1] = pairs[j-1], pairs[j]
			j--
		}
	}
	count := 0
	total := len(pairs)
	for i := 0; i < total; i++ {
		a, b := pairs[i][0], pairs[i][1]
		for j := i + 1; j < total; j++ {
			c, d := pairs[j][0], pairs[j][1]
			if a < c && c < b && b < d {
				count++
			} else if c < a && a < d && d < b {
				count++
			}
		}
	}
	return count
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
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse k: %v", idx+1, err)
		}
		if len(fields) != 2+2*k {
			return nil, fmt.Errorf("line %d: expected %d fields got %d", idx+1, 2+2*k, len(fields))
		}
		tc := testCase{n: n, k: k, pairs: make([][2]int, k)}
		pos := 2
		for i := 0; i < k; i++ {
			x, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse pair x: %v", idx+1, err)
			}
			y, err := strconv.Atoi(fields[pos+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse pair y: %v", idx+1, err)
			}
			tc.pairs[i] = [2]int{x, y}
			pos += 2
		}
		cases = append(cases, tc)
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
		expected := solveCase(tc)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for _, p := range tc.pairs {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		vals := strings.Fields(got)
		if len(vals) != 1 {
			fmt.Printf("case %d: expected single integer output, got %q\n", i+1, got)
			os.Exit(1)
		}
		gotVal, err := strconv.Atoi(vals[0])
		if err != nil {
			fmt.Printf("case %d: non-integer output %q\n", i+1, vals[0])
			os.Exit(1)
		}
		if gotVal != expected {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %d\n", i+1, expected, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
