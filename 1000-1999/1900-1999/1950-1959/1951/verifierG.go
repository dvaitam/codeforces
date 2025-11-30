package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
5 3 32 49 17 3 1
2 19 31 49
6 12 21 50 2 18 32 13
6 14 35 35 44 7 13 37
5 9 43 40 44 6 28
5 49 22 24 13 17 30
2 34 11 24
1 8 15
5 12 40 44 9 32 49
5 41 6 1 3 7 3
1 9 38
3 1 3 12
6 9 4 6 41 39 10 25
6 19 7 14 36 13 30 19
4 40 36 6 31 39
5 25 6 38 6 44 38
3 16 30 21 38
1 10 40
4 6 26 19 16 10
2 12 31 4
5 40 5 40 18 8 33
2 9 2 5
3 15 20 2 45
6 44 31 26 5 35 10
3 7 7 18 19
5 50 30 12 10 26 21
6 9 8 17 22 9 9 44
5 4 6 34 33 30 13
3 8 50 43 27
2 43 29 16
2 46 16 20
1 23 46
3 11 6 2 32
4 42 9 9 14 17
2 42 30 28
3 1 5 22 24
3 9 2 45 29
5 21 8 18 42 22 35
1 26 33
1 21 5
6 37 12 49 45 5 31 14
2 32 10 17
6 2 46 34 26 7 21 47
3 12 7 24 19
3 22 22 37 39
3 4 12 23 50
6 29 38 28 9 23 24 17
2 15 16 7
4 6 30 35 27 44
3 1 32 12 30
6 14 50 39 18 32 13 15
1 43 11
2 50 35 10
5 19 38 6 12 2 27
2 44 29 5
1 33 46
1 24 17
2 33 29 48
1 30 9
1 1 2
1 5 10
6 24 27 6 16 11 24 25
6 12 6 46 17 9 48 22
6 6 28 36 30 18 13 18
6 28 9 13 23 45 28 29
5 39 35 50 7 36 16
5 4 5 45 10 19 39
5 32 30 2 43 12 24
5 9 16 2 42 25 4
2 10 14 30
2 21 44 12
4 35 40 43 44 48
3 30 1 7 46
1 24 12
1 6 42
4 3 7 45 25 1
6 35 28 13 35 16 20 31
6 20 43 3 1 30 35 5
2 42 45 19
4 43 18 13 2 24
1 19 44
4 38 20 30 8 35
2 21 10 22
1 37 16
2 11 37 38
2 6 8 1
3 27 27 12 8
6 35 5 16 31 26 14 16
2 40 42 13
3 34 27 23 26
2 10 46 47
1 1 1
6 36 30 36 43 28 13 3
4 36 3 40 20 18
3 36 23 14 38
2 31 9 45
5 20 9 42 39 43 12
5 33 12 33 12 33 24
3 9 26 9 21
1 25 10
4 39 31 29 5 13
6 4 9 27 37 3 48 45
`

const mod int64 = 1_000_000_007

func modPow(a, e int64) int64 {
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

func modInv(a int64) int64 { return modPow(a, mod-2) }

func solve(n int, m int64, pos []int64) int64 {
	// positions may be unsorted; sort ascending
	for i := 1; i < len(pos); i++ {
		key := pos[i]
		j := i - 1
		for j >= 0 && pos[j] > key {
			pos[j+1] = pos[j]
			j--
		}
		pos[j+1] = key
	}
	var prefix, prefix2 int64
	var sum1, sum2 int64
	for j, v := range pos {
		x := v % mod
		jm := int64(j) % mod
		sum1 = (sum1 + (x*jm%mod-mod+prefix)%mod) % mod
		t := (x*x%mod*jm%mod - 2*x%mod*prefix%mod + prefix2) % mod
		if t < 0 {
			t += mod
		}
		sum2 = (sum2 + t) % mod
		prefix = (prefix + x) % mod
		prefix2 = (prefix2 + x*x%mod) % mod
	}
	total := (2 * ((m%mod)*sum1%mod - sum2 + mod) % mod) % mod
	ans := total * modInv(int64(n)%mod) % mod
	return ans
}

type testCase struct {
	n int
	m int64
	a []int64
}

func parseTestcases(raw string) ([]testCase, error) {
	fields := strings.Fields(raw)
	cases := make([]testCase, 0)
	pos := 0
	for pos < len(fields) {
		if pos+2 > len(fields) {
			return nil, fmt.Errorf("truncated data")
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("bad n")
		}
		mVal, err := strconv.ParseInt(fields[pos+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("bad m")
		}
		pos += 2
		if pos+n > len(fields) {
			return nil, fmt.Errorf("incomplete positions")
		}
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			val, err := strconv.ParseInt(fields[pos+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("bad value")
			}
			a[i] = val
		}
		pos += n
		cases = append(cases, testCase{n: n, m: mVal, a: a})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := buildInput(tc)
		expected := solve(tc.n, tc.m, append([]int64(nil), tc.a...))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strconv.FormatInt(expected, 10) {
			fmt.Printf("case %d failed: expected %d got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
