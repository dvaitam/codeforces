package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesRaw = `
3
7 4 6
8 8 5
3
2 7 8
7 6 4
4
1 10 7 1
9 6 10 8
2
7 3
1 9
3
2 1 9
1 5 10
5
6 10 9 5 6
2 8 3 6 4
2
7 1
6 9
4
3 9 8 5
10 4 2 5
2
7 7
4 7
4
2 4 4 7
3 4 1 1
10
6 10 2 6 3 7 1 10 7 10
1 9 2 9 4 10 4 2 8 7
7
1 6 8 10 5 9 8
6 2 9 6 3 5 3
8
9 5 2 9 8 9 2 9
8 1 9 9 2 10 1 3
6
8 5 6 5 5 9
6 8 8 6 4 9
7
1 1 5 6 2 6 7
3 2 6 2 4 1 2
8
1 2 1 1 3 3 1 2
2 3 1 3 2 1 2 2
10
6 9 9 1 4 8 5 5 4 10
2 5 1 10 10 2 6 8 4 8
10
4 2 5 7 7 4 8 5 3 5
2 2 8 7 3 2 6 4 10 4
8
2 5 6 7 5 8 5 8
5 9 4 3 5 9 2 1
10
9 3 7 5 8 3 6 10 5 5
1 9 9 1 7 10 8 2 10 7
2
9 9
4 8
7
10 1 10 10 6 3 3
3 3 1 8 2 3 6
5
9 1 10 3 8
3 9 8 5 10
7
7 5 7 8 8 7 1
1 2 7 2 9 6 1
3
10 3 3
5 4 3
4
10 7 2 8
1 3 9 6
3
5 4 3
2 10 6
4
5 5 3 7
9 2 7 10
10
8 10 9 4 4 5 4 6 5 2
7 6 10 10 7 8 7 8 5 3
6
8 3 1 8 8 6
8 2 5 4 1 5
2
7 10
8 2
6
10 6 4 3 9 8
5 1 6 8 5 7
3
9 1 8
4 5 10
10
9 10 5 6 6 2 2 10 8 5
10 2 5 6 2 1 6 7 4 5
7
10 4 3 5 7 1 8
5 1 10 4 8 10 10
3
5 10 8
10 1 9
5
2 10 1 5 4
10 4 7 5 9
4
6 9 8 9
2 1 6 9
6
7 1 10 1 5 9
4 10 10 10 3 5
6
8 10 8 10 10 8
8 5 7 4 5 7
5
2 8 1 1 10
6 10 9 3 7
4
8 7 1 7
10 4 6 4
9
7 8 6 7 8 8 7 5 5
1 8 3 5 5 2 7 2 7
5
8 5 7 2 3
10 9 2 5 2
8
4 4 6 10 2 1 10 8
9 3 1 8 10 1 8 5
2
9 1
8 5
9
7 10 5 10 1 1 3 7 6
9 5 10 7 2 4 5 8 10
2
3 6
4 4
2
4 1
4 7
4
2 3 6 8
1 2 2 7
7
2 9 1 5 2 5 10
9 4 2 2 6 10 4
4
6 1 9 3
8 5 4 1
5
8 7 1 4 3
8 9 6 6 3
8
8 8 6 2 10 9 2 1
8 9 10 3 1 2 7 4
9
4 4 5 2 2 6 6 6 3
9 4 6 5 3 4 9 7 7
3
2 5 8
7 6 2
8
3 5 7 3 8 6 5 9
1 5 2 5 7 2 10 4
6
4 9 5 2 9 10
9 5 5 1 3 8
10
4 10 9 2 7 2 2 4 2 5
2 2 6 8 7 5 10 10 8 10
8
5 1 3 6 1 9 8 8
1 5 9 9 8 10 5 8
8
6 10 10 7 9 5 6 7
7 5 8 6 3 9 9 4
6
8 2 6 8 3 9
3 9 3 9 4 9
6
5 7 6 10 10 10
8 8 1 3 5 3
10
6 10 9 8 10 3 7 1 6 10
9 10 1 5 8 7 6 5 3 6
8
6 2 4 7 3 4 7 8
6 1 2 10 6 1 7 10
7
3 6 6 3 3 6 5
6 2 2 2 3 8 3
9
3 3 10 9 10 4 5 5 3
4 3 2 8 6 10 2 9 7
8
9 2 3 4 9 2 7 9
3 2 9 7 1 2 10 3
4
10 8 5 7
2 5 8 10
10
1 1 5 8 6 6 5 3 7 10
9 10 8 6 4 4 5 4 7 1
`

type testCase struct {
	n int
	a []int
	b []int
}

func solveCase(tc testCase) int64 {
	total := 0
	sumSquares := 0
	for i := 0; i < tc.n; i++ {
		total += tc.a[i] + tc.b[i]
		sumSquares += tc.a[i]*tc.a[i] + tc.b[i]*tc.b[i]
	}

	dp := make([]bool, total+1)
	dp[0] = true
	for i := 0; i < tc.n; i++ {
		ndp := make([]bool, total+1)
		for s := 0; s <= total; s++ {
			if dp[s] {
				if s+tc.a[i] <= total {
					ndp[s+tc.a[i]] = true
				}
				if s+tc.b[i] <= total {
					ndp[s+tc.b[i]] = true
				}
			}
		}
		dp = ndp
	}

	best := int64(1<<63 - 1)
	for s := 0; s <= total; s++ {
		if dp[s] {
			sumA := int64(s)
			sumB := int64(total - s)
			cur := sumA*sumA + sumB*sumB
			if cur < best {
				best = cur
			}
		}
	}
	constPart := int64(sumSquares) * int64(tc.n-2)
	return best + constPart
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines)/3)
	for i := 0; i < len(lines); {
		nLine := strings.TrimSpace(lines[i])
		if nLine == "" {
			i++
			continue
		}
		n, err := strconv.Atoi(nLine)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", i+1, err)
		}
		i++
		if i >= len(lines) {
			return nil, fmt.Errorf("line %d: missing a line", i+1)
		}
		aFields := strings.Fields(lines[i])
		if len(aFields) != n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", i+1, n, len(aFields))
		}
		a := make([]int, n)
		for j, f := range aFields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a[%d]: %v", i+1, j, err)
			}
			a[j] = v
		}
		i++
		if i >= len(lines) {
			return nil, fmt.Errorf("line %d: missing b line", i+1)
		}
		bFields := strings.Fields(lines[i])
		if len(bFields) != n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", i+1, n, len(bFields))
		}
		b := make([]int, n)
		for j, f := range bFields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse b[%d]: %v", i+1, j, err)
			}
			b[j] = v
		}
		i++
		cases = append(cases, testCase{n: n, a: a, b: b})
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		for idx, v := range tc.a {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for idx, v := range tc.b {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strconv.FormatInt(expected, 10) {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
