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

const MOD = 1000000007

type queryC struct {
	tp int64
	l  int
	r  int
	x  int64
}

type caseC struct {
	n, m    int
	arr     []int64
	queries []queryC
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func fib(n int64) (int64, int64) {
	if n == 0 {
		return 0, 1
	}
	a, b := fib(n >> 1)
	c := (a * ((2*b - a + MOD) % MOD)) % MOD
	d := (a*a + b*b) % MOD
	if n&1 == 0 {
		return c, d
	}
	return d, (c + d) % MOD
}

func solve(tc caseC) []int64 {
	arr := make([]int64, tc.n)
	copy(arr, tc.arr)
	res := []int64{}
	for _, q := range tc.queries {
		if q.tp == 1 {
			for i := q.l - 1; i <= q.r-1; i++ {
				arr[i] += q.x
			}
		} else {
			sum := int64(0)
			for i := q.l - 1; i <= q.r-1; i++ {
				f, _ := fib(arr[i])
				sum = (sum + f) % MOD
			}
			res = append(res, sum)
		}
	}
	return res
}

func parseCases(path string) ([]caseC, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	cases := []caseC{}
	for {
		if !sc.Scan() {
			break
		}
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var n, m int
		fmt.Sscan(line, &n, &m)
		if !sc.Scan() {
			return nil, fmt.Errorf("bad file")
		}
		arrParts := strings.Fields(strings.TrimSpace(sc.Text()))
		if len(arrParts) < n {
			return nil, fmt.Errorf("bad array len")
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(arrParts[i], 10, 64)
			arr[i] = v
		}
		queries := make([]queryC, m)
		for i := 0; i < m; i++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("bad file")
			}
			parts := strings.Fields(strings.TrimSpace(sc.Text()))
			if len(parts) < 3 {
				return nil, fmt.Errorf("bad query")
			}
			tp, _ := strconv.ParseInt(parts[0], 10, 64)
			l, _ := strconv.Atoi(parts[1])
			r, _ := strconv.Atoi(parts[2])
			var x int64
			if tp == 1 {
				if len(parts) < 4 {
					return nil, fmt.Errorf("bad query")
				}
				x, _ = strconv.ParseInt(parts[3], 10, 64)
			}
			queries[i] = queryC{tp: tp, l: l, r: r, x: x}
		}
		cases = append(cases, caseC{n: n, m: m, arr: arr, queries: queries})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases("testcasesC.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		exp := solve(tc)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		for _, q := range tc.queries {
			if q.tp == 1 {
				fmt.Fprintf(&sb, "1 %d %d %d\n", q.l, q.r, q.x)
			} else {
				fmt.Fprintf(&sb, "2 %d %d\n", q.l, q.r)
			}
		}
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotFields := strings.Fields(out)
		if len(gotFields) != len(exp) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d values got %d\n", idx+1, len(exp), len(gotFields))
			os.Exit(1)
		}
		for i, gv := range gotFields {
			val, _ := strconv.ParseInt(gv, 10, 64)
			if val != exp[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at output %d: expected %d got %d\n", idx+1, i+1, exp[i], val)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
