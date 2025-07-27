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

type testCase struct {
	n int
	m int
	l []int
	s []int
	c []int
}

func simulateGain(cnt []int, level int, c []int) int {
	gain := c[level]
	cnt[level]++
	for lvl := level; cnt[lvl] == 2; lvl++ {
		cnt[lvl] = 0
		if lvl+1 >= len(cnt) {
			cnt = append(cnt, 0)
		}
		cnt[lvl+1]++
		gain += c[lvl+1]
	}
	return gain
}

func expected(tc testCase) int {
	cnt := make([]int, tc.n+tc.m+5)
	profit := 0
	cost := 0
	maxAllowed := tc.n + tc.m + 5
	for i := 0; i < tc.n; i++ {
		if tc.l[i] > maxAllowed {
			continue
		}
		tmp := make([]int, len(cnt))
		copy(tmp, cnt)
		g := simulateGain(tmp, tc.l[i], tc.c)
		if profit+g-(cost+tc.s[i]) > profit-cost {
			cnt = tmp
			profit += g
			cost += tc.s[i]
			if tc.l[i] < maxAllowed {
				maxAllowed = tc.l[i]
			}
		}
	}
	return profit - cost
}

func run(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.l {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.s {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i := 1; i <= tc.n+tc.m; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", tc.c[i]))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
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
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "test %d malformed\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		l := make([]int, n)
		sarr := make([]int, n)
		pos := 2
		for i := 0; i < n; i++ {
			l[i], _ = strconv.Atoi(parts[pos])
			pos++
		}
		for i := 0; i < n; i++ {
			sarr[i], _ = strconv.Atoi(parts[pos])
			pos++
		}
		c := make([]int, n+m+1)
		for i := 1; i <= n+m; i++ {
			c[i], _ = strconv.Atoi(parts[pos])
			pos++
		}
		tc := testCase{n: n, m: m, l: l, s: sarr, c: c}
		want := expected(tc)
		gotStr, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		got, _ := strconv.Atoi(gotStr)
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
