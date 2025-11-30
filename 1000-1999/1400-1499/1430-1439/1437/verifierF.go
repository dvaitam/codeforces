package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const mod = 998244353

func modinv(a, m int) int {
	b := m
	u, v := 1, 0
	for b != 0 {
		t := a / b
		a, b = b, a-t*b
		u, v = v, u-t*v
	}
	if u < 0 {
		u += m
	}
	return u
}

// solveCase mirrors 1437F.go for a single array.
func solveCase(a []int) int {
	n := len(a)
	sort.Ints(a)
	fact := make([]int, n+1)
	invf := make([]int, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = int(int64(fact[i-1]) * int64(i) % mod)
	}
	invf[n] = modinv(fact[n], mod)
	for i := n; i > 0; i-- {
		invf[i-1] = int(int64(invf[i]) * int64(i) % mod)
	}
	ans := 0
	for f := 1; f <= n; f++ {
		ok := true
		for i := f; i < n; i++ {
			if a[i] < 2*a[i-1] {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		m := n - f + 1
		F := f - 1
		c := make([]int, m+1)
		for k := 0; k < f-1; k++ {
			lo, hi := f-1, n-1
			pos := -1
			for lo <= hi {
				mid := (lo + hi) / 2
				if a[mid] >= 2*a[k] {
					pos = mid
					hi = mid - 1
				} else {
					lo = mid + 1
				}
			}
			if pos < 0 {
				ok = false
				break
			}
			t := pos - (f - 1) + 1
			if t < 1 || t > m {
				ok = false
				break
			}
			c[t]++
		}
		if !ok {
			continue
		}
		dp := make([]int, F+1)
		dp[0] = 1
		avail := 0
		for j := 1; j <= m; j++ {
			avail += c[j]
			ndp := make([]int, F+1)
			for used := 0; used <= F; used++ {
				v := dp[used]
				if v == 0 {
					continue
				}
				maxk := avail - used
				if maxk > F-used {
					maxk = F - used
				}
				for k := 0; k <= maxk; k++ {
					ways := int(int64(fact[avail-used]) * int64(invf[avail-used-k]) % mod)
					ndp[used+k] = (ndp[used+k] + v*ways) % mod
				}
			}
			dp = ndp
		}
		ans = (ans + dp[F]) % mod
	}
	return ans
}

// referenceSolve matches 1437F.go input/output behaviour.
func referenceSolve(input string) (string, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return "", fmt.Errorf("empty input")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", fmt.Errorf("bad n: %v", err)
	}
	if len(fields) != 1+n {
		return "", fmt.Errorf("expected %d tokens, got %d", 1+n, len(fields))
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[1+i])
		if err != nil {
			return "", fmt.Errorf("bad a[%d]: %v", i, err)
		}
		arr[i] = v
	}
	ans := solveCase(arr)
	return strconv.Itoa(ans), nil
}

type testCase struct {
	input string
}

// Embedded testcases from testcasesF.txt.
const testcaseData = `
3 27 6 29
3 24 8 22
1 4
5 7 2 10 14 9
1 17
5 2 1 24 5 22
2 3 14
2 23 25
5 24 20 30 28 12
1 7
2 16 25
5 30 3 28 26 17
2 28 17
4 29 6 15 20
3 8 5 4
3 16 29 17
6 23 2 29 13 7 14
4 19 5 26 25
4 2 18 26 1
5 8 18 19 12 1
2 3 17
1 8
6 23 29 6 13 12 7
1 14
2 23 25
2 23 8
4 20 14 10 12
2 24 27
6 9 10 23 29 13 12
2 30 28
2 23 7
4 22 21 30 5
4 10 25 30 27
2 6 27
1 11
6 11 29 22 20 28 30
6 2 10 7 27 16 5
6 2 23 6 18 27 7
3 15 30 14
2 17 5
2 27 22
4 17 20 27 5
6 1 20 2 13 5 2
2 22 19
4 7 16 1 25
4 17 18 4 9
4 2 27 25 11
4 28 2 16 5
6 5 23 20 12 15 11
6 2 3 7 19 22 24
6 29 1 3 14 25 9
2 4 11
2 11 26
4 18 13 18 7
2 13 27
1 16
4 27 1 15 7
1 7
2 10 20
5 6 3 16 25 20
4 13 6 1 30
1 3
1 16
5 4 3 7 13 11
3 2 25 20
6 17 19 25 22 16 17
5 24 26 27 27 22
3 21 4 22
3 25 9 1
3 2 22 28
1 13
3 1 16 19
3 30 6 8
6 30 6 17 27 20 20
2 27 9
6 15 27 15 5 7 9
2 23 4
6 8 5 3 9 15 16
3 21 1 5
2 12 23
2 2 30
1 9
2 30 19
1 15
3 12 29 19
1 21
6 13 9 11 13 11 13
3 27 11 7
5 30 10 22 23 16
6 28 7 2 2 9 28
3 28 6 20
6 19 16 29 10 30 13
2 28 30
1 26
2 30 25
5 9 6 30 28 13
3 20 22 9
1 24
3 4 19 12
6 11 27 17 23 28 11
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		if len(fields) != 1+n {
			return nil, fmt.Errorf("line %d: expected %d tokens, got %d", i+1, 1+n, len(fields))
		}
		for idx := 1; idx < len(fields); idx++ {
			if _, err := strconv.Atoi(fields[idx]); err != nil {
				return nil, fmt.Errorf("line %d: bad int at %d: %v", i+1, idx+1, err)
			}
		}
		res = append(res, testCase{input: line + "\n"})
	}
	return res, nil
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expected, err := referenceSolve(tc.input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed\ninput: %sexpected: %s\ngot: %s\n", idx+1, tc.input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
