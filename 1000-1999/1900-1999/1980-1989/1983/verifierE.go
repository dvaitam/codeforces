package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesE = `100
2 2 23 16
2 2 14 15
5 3 40 14 45 34 44
4 2 27 15 37 41
4 3 12 48 6 33
2 0 8 50
1 1 34
5 1 17 3 21 29 26
1 1 36
4 0 24 32 22 47
1 1 26
5 5 22 33 11 50 2
1 1 3
1 1 29
5 0 48 3 48 41 41
4 4 9 13 24 27
5 1 48 6 29 10 39
4 3 12 40 32 10
5 4 48 16 2 3 7
2 2 30 9
3 3 42 18 39
4 3 34 1 31 38
1 0 37
5 3 17 15 15 4 34
5 2 21 2 6 16 31
5 0 41 48 27 29 18
5 1 12 35 35 48 6
1 0 4
5 4 8 10 37 42 33
4 3 43 36 40 0
2 2 47 6
1 1 45
3 2 19 21 24
2 2 5 10
3 2 3 22 47
4 3 45 10 41 37
3 2 21 15 25
5 4 37 29 18 30 6
4 4 11 1 17 39
5 4 32 8 44 0 7
3 0 44 3 35
3 2 36 0 46
2 2 15 39
4 0 28 1 11 34
4 3 27 10 17 22
1 0 32
5 1 34 4 35 10 49
2 2 34 28
5 4 9 39 15 33 43
4 0 24 31 22 33
4 2 45 32 45 16
4 0 8 7 43 47
1 0 5
1 0 18
4 0 49 33 37 2
1 0 17
3 0 8 11 3
1 0 13
1 0 40
4 2 15 39 32 1
1 0 43
2 1 16 30
1 0 38
2 2 15 9
1 0 2
1 0 17
2 0 33 9
5 1 38 4 11 6 17
3 3 36 24 45
1 0 2
2 2 0 47
3 1 49 41 30
5 3 22 41 44 13 44
3 3 34 35 33
4 4 8 12 39 22
3 0 42 36 46
1 0 1
5 0 3 29 5 43 48
4 0 11 21 4 30
2 0 33 46
2 0 36 1
5 1 15 11 25 29 50
5 2 39 33 35 8 23
2 1 12 32
3 0 42 1 18
4 1 40 19 13 26
1 0 11
1 0 34
5 4 33 15 49 1 6
4 3 8 27 14 28
3 1 23 22 40
1 0 17
3 2 30 11 45
2 2 34 47
2 2 36 26
2 0 12 23
5 1 48 47 43 23 16
4 0 6 46 43 1
4 1 40 7 26 16
4 0 42 23 40 47
5 2 7 44 3 49 31
4 2 11 43 11 18
2 2 26 35
2 0 39 30
1 0 36`

const mod int64 = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func expected(n, k int, vals []int64) (int64, int64) {
	sum := int64(0)
	for _, v := range vals {
		sum = (sum + v) % mod
	}
	N := n - k
	invN1 := modPow(int64(N+1), mod-2)
	pSpec := (int64(N/2) + 1) % mod * invN1 % mod
	var pNorm int64
	if N > 0 {
		invN := modPow(int64(N), mod-2)
		pNorm = (int64((N+1)/2) % mod) * invN % mod
	}
	alice := int64(0)
	for i, v := range vals {
		if i < k {
			alice = (alice + (v%mod)*pSpec) % mod
		} else {
			alice = (alice + (v%mod)*pNorm) % mod
		}
	}
	bob := (sum - alice) % mod
	if bob < 0 {
		bob += mod
	}
	return alice, bob
}

type testCase struct {
	n    int
	k    int
	vals []int64
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesE)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	pos := 0
	nextInt := func() (int64, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.ParseInt(fields[pos], 10, 64)
		pos++
		return v, err
	}
	t64, err := nextInt()
	if err != nil {
		return nil, err
	}
	t := int(t64)
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		n64, err := nextInt()
		if err != nil {
			return nil, err
		}
		k64, err := nextInt()
		if err != nil {
			return nil, err
		}
		n := int(n64)
		k := int(k64)
		vals := make([]int64, n)
		for j := 0; j < n; j++ {
			v, err := nextInt()
			if err != nil {
				return nil, err
			}
			vals[j] = v
		}
		tests[i] = testCase{n: n, k: k, vals: vals}
	}
	return tests, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d", tc.n, tc.k)
		for _, v := range tc.vals {
			sb.WriteByte(' ')
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	tests, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}
	input := buildInput(tests)
	output, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outFields := strings.Fields(output)
	if len(outFields) != 2*len(tests) {
		fmt.Printf("expected %d numbers, got %d\n", 2*len(tests), len(outFields))
		os.Exit(1)
	}
	idx := 0
	for i, tc := range tests {
		al, bo := expected(tc.n, tc.k, tc.vals)
		if idx+1 >= len(outFields) {
			fmt.Printf("missing output for case %d\n", i+1)
			os.Exit(1)
		}
		gotA, err := strconv.ParseInt(outFields[idx], 10, 64)
		if err != nil {
			fmt.Printf("invalid output for case %d alice: %v\n", i+1, err)
			os.Exit(1)
		}
		gotB, err := strconv.ParseInt(outFields[idx+1], 10, 64)
		if err != nil {
			fmt.Printf("invalid output for case %d bob: %v\n", i+1, err)
			os.Exit(1)
		}
		if gotA != al || gotB != bo {
			fmt.Printf("case %d failed\nexpected: %d %d\ngot: %d %d\n", i+1, al, bo, gotA, gotB)
			os.Exit(1)
		}
		idx += 2
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
