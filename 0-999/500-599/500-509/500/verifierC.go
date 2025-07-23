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

type testCaseC struct {
	n int
	m int
	w []int64
	q []int
}

func parseTestcases(path string) ([]testCaseC, error) {
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
	cases := make([]testCaseC, T)
	for i := 0; i < T; i++ {
		var n, m int
		if _, err := fmt.Fscan(in, &n, &m); err != nil {
			return nil, err
		}
		w := make([]int64, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &w[j])
		}
		q := make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &q[j])
		}
		cases[i] = testCaseC{n: n, m: m, w: w, q: q}
	}
	return cases, nil
}

func solveCase(tc testCaseC) string {
	n := tc.n
	last := make([]int, n+1)
	var ans int64
	for t, x := range tc.q {
		for i := 1; i <= n; i++ {
			if i != x && last[i] > last[x] {
				ans += tc.w[i-1]
			}
		}
		last[x] = t + 1
	}
	return strconv.FormatInt(ans, 10)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for j, v := range tc.w {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for j, v := range tc.q {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		expected := solveCase(tc)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
