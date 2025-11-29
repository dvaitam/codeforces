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

const (
	MOD          int64 = 998244353
	testcasesRaw       = `100
2
8 5
1
1
3
10 8 6
6
1 5 8 4 7 9
2
4 10
5
10 2 7 6 2
6
7 5 8 2 4 5
2
1 10
4
6 8 4 9
1
6
4
10 7 5 6
2
2 9
4
2 10 5 5
4
7 8 4 3
4
9 1 4 3
1
6
5
6 7 9 7 5
3
8 1 3
7
10 7 2 8 4 2 10
8
8 7 2 9 7 8 5 7
2
4 5
8
8 3 1 1 9 2 5 10
6
4 5 9 8 6 9
5
7 7 10 8 5
8
8 8 3 7 8 5 8 6
6
3 10 7 10 5 6
7
8 3 5 9 1 10 8
1
3
1
10
2
6 6
8
10 1 4 3 5 10 1 7
8
2 8 4 2 6 6 3 4
6
3 1 2 2 1 8
8
2 1 9 3 10 2 3 10
5
8 8 1 3 4
8
8 8 8 9 5 8 10 2
5
6 5 6 1 2
5
5 5 8 7 8
6
1 2 4 10 4 9
1
3
6
6 1 1 4 7 6
2
3 2
7
9 3 8 8 2 3 2
6
9 1 3 9 3 3
5
3 9 10 7 4
4
10 3 7 3
6
4 9 10 1 10 7
4
6 7 9 8
1
7
3
4 9 6
5
3 3 3 3 7
2
1 1
1
8
1
2
2
6 10
1
1
4
7 4 10 5
2
8 2
1
3
8
2 2 3 1 1 7 5 5
8
4 4 9 4 5 10 1 5
6
9 2 10 10 1 6
6
2 4 1 6 6 7
1
6
3
1 1 6
2
8 1
8
5 2 10 10 5 4 3 3
7
9 3 4 7 10 7 10
8
5 5 3 5 1 4 8 2
3
10 2 1
3
2 3 10
7
8 8 9 9 8 4 7
3
8 1 6
7
5 6 3 5 6 10 6
4
1 5 3 2
7
10 10 4 2 4 7 7
3
3 9 1
8
6 6 10 7 5 2 5 3
1
2
3
9 9 2
2
2 2
4
6 6 3 4
3
3 4 3
8
8 10 2 10 3 8 4 10
6
3 6 5 2 5 1
1
1
8
8 6 1 1 5 2 7 10
6
2 6 3 3 6 8
7
1 9 3 2 8 6 1
3
9 8 5
2
8 2
7
4 9 3 8 6 5 9
6
5 9 5 8 9 1
4
10 1 9 8
5
3 1 6 1 4
5
4 3 5 1 5
5
6 4 3 10 1
2
6 10
3
6 3 8
4
8 7 9 10`
)

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

func parseTests() ([]testCaseG, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("missing test count")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	cases := make([]testCaseG, 0, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d: missing n", i+1)
		}
		n, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %v", i+1, err)
		}
		vals := make([]int64, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d: missing value %d", i+1, j+1)
			}
			v, err := strconv.ParseInt(scan.Text(), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d: bad value %d: %v", i+1, j+1, err)
			}
			vals[j] = v
		}
		cases = append(cases, testCaseG{vals: vals})
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

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
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

	cases, err := parseTests()
	if err != nil {
		fmt.Println("failed to parse embedded tests:", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", len(tc.vals)))
		for i, v := range tc.vals {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", v))
		}
		input.WriteByte('\n')

		expected := solveCase(tc)
		outStr, err := runCandidate(bin, []byte(input.String()))
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(outStr)
		if len(fields) != len(expected) {
			fmt.Printf("case %d: expected %d numbers got %d\n", idx+1, len(expected), len(fields))
			os.Exit(1)
		}
		for i, f := range fields {
			v, err := strconv.ParseInt(f, 10, 64)
			if err != nil || v != expected[i] {
				fmt.Printf("case %d index %d expected %d got %s\n", idx+1, i+1, expected[i], f)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
