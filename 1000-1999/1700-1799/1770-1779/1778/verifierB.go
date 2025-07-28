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

type testCaseB struct {
	n, m, d int
	p       []int
	a       []int
}

func parseCasesB(path string) ([]testCaseB, error) {
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
	cases := make([]testCaseB, t)
	for i := 0; i < t; i++ {
		var n, m, d int
		if _, err := fmt.Fscan(in, &n, &m, &d); err != nil {
			return nil, err
		}
		p := make([]int, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(in, &p[j]); err != nil {
				return nil, err
			}
		}
		a := make([]int, m)
		for j := 0; j < m; j++ {
			if _, err := fmt.Fscan(in, &a[j]); err != nil {
				return nil, err
			}
		}
		cases[i] = testCaseB{n: n, m: m, d: d, p: p, a: a}
	}
	return cases, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveB(tc testCaseB) int {
	pos := make([]int, tc.n+1)
	for i, v := range tc.p {
		pos[v] = i + 1
	}
	ans := tc.n + 5
	for i := 0; i < tc.m-1; i++ {
		x := pos[tc.a[i]]
		y := pos[tc.a[i+1]]
		if x > y || y-x > tc.d {
			return 0
		}
		diff := y - x
		cur := diff
		delta := tc.d + 1 - diff
		if y+delta <= tc.n {
			cur = min(cur, delta)
		}
		ans = min(ans, cur)
	}
	if ans > tc.n {
		ans = 0
	}
	return ans
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCasesB("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		sb := strings.Builder{}
		fmt.Fprintf(&sb, "1\n%d %d %d\n", tc.n, tc.m, tc.d)
		for i, v := range tc.p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveB(tc)
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(gotStr))
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
