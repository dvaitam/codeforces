package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesRaw = `4 6 0 4 8
4 6 4 7 5
5 3 8 2 4 2
1 9
3 8 9 2
3 1 1 10
3 7 8 1
3 6 5 9
2 8 7
4 8 4 0 8
1 1
4 10 10 0 9
4 5 3 5 1
2 9 3
2 2 8
4 1 1 5 8
4 1 4 8 4
1 8
3 8 3 9
5 9 4 7 1 9
4 9 8 2 0
3 2 3 2
3 5 10 10
1 10
3 1 10 9
1 8
5 3 3 8 9 10
2 0 2
4 8 2 4 0
2 7 9
5 8 1 2 2 8
2 3 8
5 0 6 2 4 4
2 3 5
5 0 7 8 1 1
2 6 10
3 8 7 1
4 9 7 2 5
5 1 0 1 10 2
4 6 7 9 5
3 1 2 9
5 3 0 6 6 0
1 1
5 0 6 2 0 4
1 7
5 8 5 1 10 6
2 5 3
5 4 4 4 8 0
5 6 8 9 10 1
3 1 6 9
5 10 1 6 1 2
3 9 6 9
4 6 9 9 1
1 0
3 5 10 10
3 0 7 8
2 8 6
5 9 7 7 1 8
3 3 1 4
3 0 6 2
1 0
4 6 10 3 2
1 6
1 6
1 1
5 5 1 8 10 8
3 9 3 6
2 3 9
3 2 8 7
4 0 10 1 1
1 10
4 2 2 10 5
4 5 10 0 2
1 9
4 6 7 4 6
2 2 10
3 4 5 7
5 7 6 0 10 8
1 0
5 9 4 4 0 1
2 6 2
5 9 9 8 0 0
5 2 4 1 2 10
1 3
2 10 5
2 3 4
5 7 7 8 5 0
5 0 1 1 3 0
3 8 6 7
4 6 6 3 0
5 1 4 7 2 1
2 9 0
3 1 8 6
2 4 6
5 4 9 7 0 9
2 2 2
1 2
5 6 6 5 9 10
3 3 10 7
1 0
2 4 4
1 10
2 2 10
3 8 1 9
3 9 2 5
1 2
3 7 6 4
3 10 1 9
2 9 5
3 8 10 8
4 5 2 1 8
3 0 9 9
4 9 0 3 4
3 5 6 10
1 6
5 4 4 10 7 7`

type testCase struct {
	n   int
	arr []int
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveCase(tc testCase) []int {
	const maxBits = 30
	counts := make([]int, maxBits)
	for _, x := range tc.arr {
		for b := 0; b < maxBits; b++ {
			if (x>>b)&1 == 1 {
				counts[b]++
			}
		}
	}
	g := 0
	for _, c := range counts {
		g = gcd(g, c)
	}
	if g == 0 {
		res := make([]int, tc.n)
		for k := 1; k <= tc.n; k++ {
			res[k-1] = k
		}
		return res
	}
	res := make([]int, 0, tc.n)
	for k := 1; k <= tc.n; k++ {
		if g%k == 0 {
			res = append(res, k)
		}
	}
	return res
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, n, len(fields)-1)
		}
		tc := testCase{n: n, arr: make([]int, n)}
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse value: %v", idx+1, err)
			}
			tc.arr[i] = v
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
		for idx, v := range tc.arr {
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
		vals := strings.Fields(got)
		if len(vals) != len(expected) {
			fmt.Printf("case %d: expected %d outputs, got %d\n", i+1, len(expected), len(vals))
			os.Exit(1)
		}
		for idx, v := range vals {
			num, err := strconv.Atoi(v)
			if err != nil || num != expected[idx] {
				fmt.Printf("case %d mismatch at position %d: expected %d got %s\n", i+1, idx+1, expected[idx], v)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
