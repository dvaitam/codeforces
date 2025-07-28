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

type testCaseF struct {
	n int
	c int64
	a []int64
	b []int64
}

func parseTestcasesF(path string) ([]testCaseF, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cases []testCaseF
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, fmt.Errorf("bad line: %s", line)
		}
		n, _ := strconv.Atoi(fields[0])
		cVal, _ := strconv.ParseInt(fields[1], 10, 64)
		if len(fields)-2 != 2*n-1 {
			return nil, fmt.Errorf("bad count")
		}
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i], _ = strconv.ParseInt(fields[2+i], 10, 64)
		}
		b := make([]int64, n-1)
		for i := 0; i < n-1; i++ {
			b[i], _ = strconv.ParseInt(fields[2+n+i], 10, 64)
		}
		cases = append(cases, testCaseF{n: n, c: cVal, a: a, b: b})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func solveF(n int, c int64, a, b []int64) string {
	curDay := int64(0)
	curMoney := int64(0)
	ans := int64(1 << 60)
	for i := 0; i < n; i++ {
		if curMoney >= c {
			if curDay < ans {
				ans = curDay
			}
		} else {
			need := c - curMoney
			days := curDay + (need+a[i]-1)/a[i]
			if days < ans {
				ans = days
			}
		}
		if i == n-1 {
			break
		}
		if curMoney < b[i] {
			need := b[i] - curMoney
			d := (need + a[i] - 1) / a[i]
			curDay += d
			curMoney += d * a[i]
		}
		curMoney -= b[i]
		curDay++
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcasesF("testcasesF.txt")
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.c))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		expected := solveF(tc.n, tc.c, tc.a, tc.b)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
