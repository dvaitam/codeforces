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
const testcasesRaw = `2 2
2 2
1 2
2 2
2 4
1 2
1 2
1 3
2 2
3 3
1 3
1 2
3 2
3 2
1 2
2 5
1 2
2 5
3 3
3 2
1 4
3 5
3 3
1 2
3 3
1 5
2 4
3 5
2 3
3 2
1 5
2 4
2 4
1 2
1 3
2 3
3 3
1 2
2 5
3 3
2 4
1 2
1 4
3 4
1 2
3 5
1 4
2 3
3 2
3 3
2 5
1 4
2 4
1 4
1 5
1 5
3 4
1 4
1 4
3 2
2 5
1 4
1 3
2 4
3 2
1 5
1 2
2 4
3 2
2 2
1 3
2 3
2 3
3 3
1 3
2 5
3 4
2 4
1 2
2 4
1 5
1 5
3 2
3 3
3 5
2 5
1 3
3 5
2 5
3 3
2 2
3 3
3 2
2 3
1 3
2 3
2 4
3 5
2 5
2 5`

const mod int64 = 1000000007

type testCase struct {
	n int
	m int
}

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func solveCase(n, m int) int64 {
	if m < 2 {
		return 0
	}
	maxN := m + 2*n + 5
	fac := make([]int64, maxN+1)
	ifac := make([]int64, maxN+1)
	fac[0] = 1
	for i := 1; i <= maxN; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	ifac[maxN] = modPow(fac[maxN], mod-2)
	for i := maxN - 1; i >= 0; i-- {
		ifac[i] = ifac[i+1] * int64(i+1) % mod
	}
	comb := func(nv, kv int) int64 {
		if nv < 0 || kv < 0 || kv > nv {
			return 0
		}
		return fac[nv] * ifac[kv] % mod * ifac[nv-kv] % mod
	}

	A := make([]int64, n)
	Astar := make([]int64, n)
	S0 := make([]int64, n)
	S1 := make([]int64, n)

	var ans int64
	for W := 0; W <= m-2; W++ {
		for k := 0; k < n; k++ {
			top := W + 2*k
			A[k] = comb(top, 2*k)
		}
		for k := 0; k < n; k++ {
			if k == 0 {
				Astar[k] = 1
			} else {
				t1 := comb(W-1+2*k, 2*k)
				t2 := comb(W-2+2*k, 2*k)
				val := (2*t1 - t2) % mod
				if val < 0 {
					val += mod
				}
				Astar[k] = val
			}
		}
		var s0, s1 int64
		for i := 0; i < n; i++ {
			s0 += A[i]
			if s0 >= mod {
				s0 -= mod
			}
			S0[i] = s0
			s1 = (s1 + int64(i)*A[i]) % mod
			S1[i] = s1
		}
		var T int64
		for k := 0; k < n; k++ {
			Nk := n - 1 - k
			rk := (int64(n-k)*S0[Nk] - S1[Nk]) % mod
			if rk < 0 {
				rk += mod
			}
			T = (T + Astar[k]*rk) % mod
		}
		mult := int64(m - W - 1)
		ans = (ans + T*mult) % mod
	}
	if ans < 0 {
		ans += mod
	}
	return ans % mod
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d malformed", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %w", idx+1, err)
		}
		cases = append(cases, testCase{n: n, m: m})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d %d\n", tc.n, tc.m)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range testcases {
		expect := strconv.FormatInt(solveCase(tc.n, tc.m), 10)
		input := buildInput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\nInput:\n%s\n", idx+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
