package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const MOD int64 = 998244353

func powMod(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

type testCaseG struct {
	vals []int64
}

func parseTestcases(path string) ([]testCaseG, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return nil, err
	}
	cases := make([]testCaseG, T)
	for i := 0; i < T; i++ {
		var n int
		fmt.Fscan(in, &n)
		vals := make([]int64, n)
		for j := 0; j < n; j++ {
			var x int
			fmt.Fscan(in, &x)
			vals[j] = int64(x)
		}
		cases[i] = testCaseG{vals: vals}
	}
	return cases, nil
}

func solveCase(tc testCaseG) []int64 {
	vals := append([]int64(nil), tc.vals...)
	n := len(vals)
	sort.Slice(vals, func(i, j int) bool { return vals[i] > vals[j] })
	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + vals[i-1]
	}
	invN := powMod(int64(n), MOD-2)
	res := make([]int64, n)
	for k := 1; k <= n; k++ {
		q := (n - k) / k
		r := (n - k) % k
		var sum int64
		for t := 0; t < q; t++ {
			l := k + t*k + 1
			rIndex := k + (t+1)*k
			segment := prefix[rIndex] - prefix[l-1]
			sum += int64(t+1) * segment
		}
		if r > 0 {
			l := k + q*k + 1
			rIndex := k + q*k + r
			segment := prefix[rIndex] - prefix[l-1]
			sum += int64(q+1) * segment
		}
		res[k-1] = sum % MOD * invN % MOD
	}
	return res
}

func run(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return out.String(), errBuf.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases("testcasesG.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.vals)))
		for i, v := range tc.vals {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		outStr, errStr, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\n", idx+1, err, errStr)
			os.Exit(1)
		}
		expected := solveCase(tc)
		fields := strings.Fields(strings.TrimSpace(outStr))
		if len(fields) != len(expected) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d numbers got %d\n", idx+1, len(expected), len(fields))
			os.Exit(1)
		}
		for i, f := range fields {
			v, err := strconv.ParseInt(f, 10, 64)
			if err != nil || v != expected[i] {
				fmt.Fprintf(os.Stderr, "case %d index %d expected %d got %s\n", idx+1, i+1, expected[i], f)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
