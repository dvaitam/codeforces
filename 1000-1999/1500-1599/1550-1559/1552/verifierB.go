package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n        int
	athletes [][5]int
}

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `
2 2 1 2 1 1 2 1 2 1 1
1 1 1 1 1 1
5 1 1 1 3 2 2 5 5 5 4 1 3 2 3 1 1 2 4 5 4 1 1 4 5 5
1 1 1 1 1 1
5 3 2 5 4 2 1 5 3 1 5 5 5 5 5 1 4 4 1 4 3 3 1 3 1 3
2 1 2 1 1 2 2 1 2 2 1
2 1 2 1 2 1 1 1 1 2 2
5 4 5 4 3 1 2 3 3 1 3 3 1 3 4 4 4 4 3 2 4 4 3 5 3 1
4 1 4 2 3 3 1 1 3 3 3 4 4 2 4 3 4 1 3 4 3
1 1 1 1 1 1
2 2 1 1 1 1 2 2 2 1 1
4 3 3 4 1 2 3 4 3 3 2 1 1 2 4 4 1 3 3 4 2
5 3 4 2 1 2 1 4 2 4 1 3 3 1 5 5 3 3 2 2 3 5 2 1 2 2
2 2 1 2 2 2 1 1 1 2 2
4 4 4 1 2 3 1 4 4 3 2 2 2 4 1 1 2 3 2 2 3
2 2 2 1 2 1 2 1 2 1 1
5 5 2 2 5 5 3 5 1 1 1 1 5 3 2 5 4 5 1 4 5 3 3 4 2 4
3 2 1 3 1 1 2 2 2 2 1 2 3 1 2 3
3 3 2 2 1 2 3 3 1 2 1 2 1 2 2 1
2 1 1 2 2 2 1 1 2 2 1
2 2 2 2 2 2 1 1 2 2 1
2 2 1 1 2 2 2 2 1 1 2
2 2 2 2 2 1 2 1 1 2 2
1 1 1 1 1 1
1 1 1 1 1 1
2 2 1 1 1 1 2 2 2 1 1
3 3 2 1 2 3 3 2 3 2 3 3 3 2 3 3
3 2 1 2 2 3 1 2 3 2 3 3 3 1 3 2
3 3 3 3 1 3 1 3 2 1 3 3 3 3 2 1
2 1 1 1 1 2 2 2 1 2 1
4 4 1 2 2 3 4 4 3 2 4 3 3 3 4 2 2 3 3 2 3
3 3 1 3 3 2 2 2 2 3 3 2 2 1 3 3
5 4 4 3 2 3 1 3 1 3 4 3 4 2 2 5 3 5 3 2 1 5 5 2 2 1
5 2 2 1 1 4 3 4 1 2 1 5 2 4 4 4 2 3 3 3 1 1 5 2 5 1
2 2 1 1 2 2 1 2 1 1 2
5 5 1 3 3 4 5 5 4 2 5 4 3 2 1 3 2 4 3 2 1 3 3 5 2 2
5 4 4 5 3 1 4 1 2 2 4 3 2 1 3 4 2 4 3 3 3 5 5 1 3 3
2 1 1 2 2 1 2 2 1 2 2
5 2 2 2 4 2 4 5 3 4 5 4 5 2 2 5 1 5 4 5 3 1 5 1 1 5
5 2 5 5 2 2 3 5 4 1 2 5 4 1 5 4 1 4 2 5 4 4 5 1 5 4
3 3 1 2 2 1 3 1 3 1 1 2 2 3 1 3
1 1 1 1 1 1
2 1 2 2 2 2 2 2 1 1 2
1 1 1 1 1 1
1 1 1 1 1 1
2 2 2 1 2 1 1 1 2 1 1
2 1 1 1 2 1 1 1 2 2 1
5 2 1 2 5 3 2 5 3 5 5 3 2 3 2 5 2 4 3 3 1 5 5 2 4 5
4 1 3 4 3 4 4 2 3 4 2 4 2 2 3 2 4 1 4 4 2
4 2 3 2 3 1 4 3 1 2 1 3 4 1 1 2 2 4 2 3 3
1 1 1 1 1 1
5 2 4 2 5 3 3 4 4 5 3 3 3 3 5 4 4 1 4 3 2 3 3 4 1 5
2 1 1 2 2 2 1 1 2 1 2
5 2 4 5 3 3 4 4 4 2 3 2 2 5 4 2 5 1 5 1 1 1 1 4 2 2
1 1 1 1 1 1
1 1 1 1 1 1
1 1 1 1 1 1
2 1 1 1 1 1 2 2 2 2 1
3 2 2 3 1 3 1 2 3 1 3 3 1 1 2 2
4 4 1 1 2 1 4 4 3 1 1 2 2 4 3 1 1 3 1 2 3
1 1 1 1 1 1
2 1 2 2 2 1 1 2 1 1 2
2 1 2 2 2 1 2 2 2 1 2
5 4 4 2 5 4 1 5 3 2 1 1 3 2 4 2 4 3 3 1 1 1 3 4 3 3
1 1 1 1 1 1
5 5 5 4 1 1 1 3 5 1 5 5 3 4 1 3 4 4 1 5 2 5 5 2 4 2
2 2 2 2 2 1 1 1 2 2 2
3 2 1 2 2 1 3 3 2 3 3 1 2 1 1 2
1 1 1 1 1 1
3 2 3 3 1 2 2 3 3 1 2 2 1 2 2 2
4 2 3 4 3 4 3 1 2 2 1 4 2 4 3 1 4 4 2 3 2
3 3 1 3 2 2 1 2 3 2 3 2 3 3 3 2
5 1 3 5 5 1 1 2 3 2 5 1 5 4 3 4 1 1 1 3 2 5 5 5 5 5
2 2 2 2 2 1 2 1 2 1 2
4 1 1 2 3 3 3 4 2 4 3 4 1 4 4 3 1 3 1 1 1
3 3 1 3 1 3 1 3 1 3 1 2 3 2 2 1
4 2 2 4 4 1 2 4 4 4 1 2 2 3 1 1 2 3 3 2 4
1 1 1 1 1 1
5 1 3 3 4 5 2 3 5 5 1 2 3 1 2 3 1 2 3 2 5 4 4 5 2 5
4 4 2 4 2 2 4 3 2 2 3 4 1 3 1 4 2 2 4 4 4
4 4 4 4 2 1 1 2 3 1 3 2 3 2 4 4 1 1 3 3 1
4 2 4 1 2 3 1 4 3 3 1 4 1 2 2 3 4 2 2 2 2
1 1 1 1 1 1
3 3 2 3 3 1 1 2 1 2 3 1 3 1 1 1
3 1 3 2 2 1 3 3 2 3 2 2 2 1 3 1
4 3 3 4 4 3 1 1 4 3 1 4 2 4 2 4 3 3 3 3 2
5 4 3 1 1 1 4 1 1 2 4 4 1 4 2 2 3 2 3 1 3 3 1 3 2 1
4 3 2 4 1 1 1 2 4 1 2 2 4 1 1 3 1 4 2 4 1
2 2 1 1 1 2 2 2 1 2 1
2 1 1 1 2 2 1 1 2 2 1
1 1 1 1 1 1
4 3 2 2 3 1 2 3 1 3 2 2 4 3 1 2 4 4 1 2 4
5 1 5 1 3 3 2 4 5 4 3 5 2 4 1 1 5 1 1 3 2 1 1 4 5 3
5 4 4 4 1 5 5 5 2 3 1 3 1 5 2 2 5 3 4 5 1 3 4 2 1 2
1 1 1 1 1 1
2 1 1 1 1 1 2 2 2 1 2
2 2 1 1 2 2 2 1 2 1 2
4 3 1 4 2 3 2 3 1 1 2 1 1 1 2 1 4 3 4 2 2
1 1 1 1 1 1
3 2 1 3 2 1 3 3 1 1 1 3 1 2 1 2

`

func better(a, b [5]int) bool {
	cnt := 0
	for i := 0; i < 5; i++ {
		if a[i] < b[i] {
			cnt++
		}
	}
	return cnt >= 3
}

func solveCase(tc testCase) int {
	candidate := 0
	for i := 1; i < tc.n; i++ {
		if better(tc.athletes[i], tc.athletes[candidate]) {
			candidate = i
		}
	}
	for i := 0; i < tc.n; i++ {
		if i == candidate {
			continue
		}
		if !better(tc.athletes[candidate], tc.athletes[i]) {
			return -1
		}
	}
	return candidate + 1
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
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		if len(fields) != 1+n*5 {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, 1+n*5, len(fields))
		}
		tc := testCase{n: n, athletes: make([][5]int, n)}
		pos := 1
		for i := 0; i < n; i++ {
			for j := 0; j < 5; j++ {
				v, err := strconv.Atoi(fields[pos])
				if err != nil {
					return nil, fmt.Errorf("line %d: parse val: %v", idx+1, err)
				}
				tc.athletes[i][j] = v
				pos++
			}
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for idx, ath := range tc.athletes {
			for j := 0; j < 5; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(ath[j]))
			}
			if idx+1 == tc.n {
				sb.WriteByte('\n')
			} else {
				sb.WriteByte('\n')
			}
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
