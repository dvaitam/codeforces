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
const testcasesRaw = `100
5 3
5 3
3 4
2 4
5 4
3 2
2 5
5 4
1 4
3 2
5 2
2 2
4 4
4 3
2 1
3 1
5 4
3 1
5 4
5 1
1 1
1 4
2 3
1 5
5 4
2 2
1 4
4 2
2 1
5 1
4 2
2 2
3 3
1 5
1 1
4 5
4 1
1 5
4 4
4 3
5 1
2 3
5 1
2 1
3 2
2 1
5 2
4 4
1 3
3 1
4 5
3 3
2 1
2 4
3 4
2 5
5 3
2 3
5 5
3 3
3 5
2 2
2 1
5 5
3 4
2 5
2 1
1 4
2 4
4 1
4 1
3 4
4 5
3 4
1 3
2 4
5 1
1 4
2 5
5 2
5 1
5 1
2 5
1 1
5 1
1 3
4 5
3 3
3 3
2 3
2 2
4 5
1 1
5 3
2 1
5 3
2 5
3 2
1 4
2 1`

const mod int64 = 998244353

type testCase struct {
	n int64
	m int64
}

func comb(n, k int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	res := int64(1)
	for i := int64(1); i <= k; i++ {
		res = res * (n - i + 1) / i
	}
	return res % mod
}

func solveCase(n, m int64) int64 {
	if n == 1 || m == 1 {
		return 0
	}
	if n == 2 {
		return 2 * comb(m+2, 4) % mod
	}
	if m == 2 {
		return 2 * comb(n+2, 4) % mod
	}
	return 0
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid t: %v", err)
	}
	res := make([]testCase, 0, t)
	idx := 1
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if idx+1 >= len(fields) {
			return nil, fmt.Errorf("case %d: missing numbers", caseIdx+1)
		}
		n, err := strconv.ParseInt(fields[idx], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", caseIdx+1, err)
		}
		m, err := strconv.ParseInt(fields[idx+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d: parse m: %v", caseIdx+1, err)
		}
		idx += 2
		res = append(res, testCase{n: n, m: m})
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("extra data after parsing")
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
		expected := solveCase(tc.n, tc.m)
		input := fmt.Sprintf("%d %d\n", tc.n, tc.m)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if err != nil || val != expected {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
