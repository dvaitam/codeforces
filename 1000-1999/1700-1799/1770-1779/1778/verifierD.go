package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD int64 = 998244353

type testCaseD struct {
	n    int
	a, b string
}

func parseCasesD(path string) ([]testCaseD, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, err
	}
	cases := make([]testCaseD, t)
	for i := 0; i < t; i++ {
		var n int
		var s1, s2 string
		if _, err := fmt.Fscan(in, &n); err != nil {
			return nil, err
		}
		if _, err := fmt.Fscan(in, &s1); err != nil {
			return nil, err
		}
		if _, err := fmt.Fscan(in, &s2); err != nil {
			return nil, err
		}
		cases[i] = testCaseD{n: n, a: s1, b: s2}
	}
	return cases, nil
}

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func expectedMoves(n, d int) int64 {
	if d == 0 {
		return 0
	}
	inv := make([]int64, n+1)
	inv[1] = 1
	for i := 2; i <= n; i++ {
		inv[i] = MOD - (MOD/int64(i))*inv[MOD%int64(i)]%MOD
	}
	A := make([]int64, n+1)
	B := make([]int64, n+1)
	if n >= 1 {
		A[1] = 1
	}
	for i := 1; i < n; i++ {
		denom := n - i
		x1 := (int64(n)*A[i] - int64(i)*A[i-1]) % MOD
		if x1 < 0 {
			x1 += MOD
		}
		A[i+1] = x1 * inv[denom] % MOD

		x2 := (int64(n)*B[i] - int64(i)*B[i-1] - int64(n)) % MOD
		if x2 < 0 {
			x2 += MOD
		}
		B[i+1] = x2 * inv[denom] % MOD
	}
	diff := (A[n] - A[n-1]) % MOD
	if diff < 0 {
		diff += MOD
	}
	x := (B[n-1] + 1 - B[n]) % MOD
	if x < 0 {
		x += MOD
	}
	x = x * modPow(diff, MOD-2) % MOD
	res := (A[d]*x + B[d]) % MOD
	return res
}

func solveD(tc testCaseD) int64 {
	d := 0
	for i := 0; i < tc.n; i++ {
		if tc.a[i] != tc.b[i] {
			d++
		}
	}
	return expectedMoves(tc.n, d)
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCasesD("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		sb := strings.Builder{}
		fmt.Fprintf(&sb, "1\n%d\n%s\n%s\n", tc.n, tc.a, tc.b)
		input := sb.String()
		expected := solveD(tc)
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(gotStr), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
