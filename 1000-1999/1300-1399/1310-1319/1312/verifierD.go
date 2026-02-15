package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod = 998244353

const testcasesRaw = `5 9
7 16
9 19
3 3
9 13
10 13
5 20
10 18
9 15
4 11
4 20
8 19
2 4
4 5
6 18
2 10
9 18
8 19
8 14
9 11
2 20
7 13
9 15
3 10
10 15
9 16
6 16
10 12
2 17
8 10
2 6
4 8
9 10
4 14
10 14
5 16
8 14
5 17
4 13
4 3
8 11
7 12
4 6
5 13
9 15
9 16
3 11
8 19
7 17
9 15
7 12
9 18
8 12
10 19
5 7
10 8
3 14
3 9
2 5
7 13
2 5
2 7
9 20
10 6
4 20
5 13
3 8
9 12
10 3
6 4
8 4
10 3
7 12
5 5
7 5
6 9
5 14
4 14
2 9
7 17
3 15
8 6
2 9
3 13
3 19
8 5
6 5
9 5
5 11
8 9
7 4
2 16
4 15
9 19
7 9
8 9
9 12
9 8
10 7
8 7`

// Embedded reference logic from 1312D.go.
func modpow(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func solve(n, m int) int64 {
	if n < 3 {
		return 0
	}
	fac := make([]int64, m+1)
	ifac := make([]int64, m+1)
	fac[0] = 1
	for i := 1; i <= m; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	ifac[m] = modpow(fac[m], mod-2)
	for i := m; i > 0; i-- {
		ifac[i-1] = ifac[i] * int64(i) % mod
	}

	// Formula: nCr(m, n-1) * (n-2) * 2^(n-3)
	if n-1 > m {
		return 0
	}
	cmb := fac[m] * ifac[n-1] % mod * ifac[m-(n-1)] % mod
	term2 := int64(n - 2)
	pow2 := modpow(2, int64(n-3))

	return cmb * term2 % mod * pow2 % mod
}

type testCase struct {
	n int
	m int
}

func parseTestcases(raw string) ([]testCase, error) {
	fields := strings.Fields(raw)
	if len(fields)%2 != 0 {
		return nil, fmt.Errorf("expected pairs of ints, got %d tokens", len(fields))
	}
	res := make([]testCase, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		n, err1 := strconv.Atoi(fields[i])
		m, err2 := strconv.Atoi(fields[i+1])
		if err := firstErr(err1, err2); err != nil {
			return nil, fmt.Errorf("invalid numbers at position %d: %w", i+1, err)
		}
		res = append(res, testCase{n: n, m: m})
	}
	return res, nil
}

func firstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test data:", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expected := fmt.Sprintf("%d", solve(tc.n, tc.m))
		input := fmt.Sprintf("%d %d\n", tc.n, tc.m)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
