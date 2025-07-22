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

func solveQueries(a []int, queries []int) []int {
	n := len(a)
	ans := make([]int, n+1)
	seen := make(map[int]bool)
	distinct := 0
	for i := n - 1; i >= 0; i-- {
		if !seen[a[i]] {
			seen[a[i]] = true
			distinct++
		}
		ans[i] = distinct
	}
	res := make([]int, len(queries))
	for i, q := range queries {
		res[i] = ans[q-1]
	}
	return res
}

func runCase(bin, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read testcasesB.txt: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())

	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
			os.Exit(1)
		}
		m, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
				os.Exit(1)
			}
			arr[i], _ = strconv.Atoi(scan.Text())
		}
		queries := make([]int, m)
		for i := 0; i < m; i++ {
			if !scan.Scan() {
				fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
				os.Exit(1)
			}
			queries[i], _ = strconv.Atoi(scan.Text())
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, q := range queries {
			fmt.Fprintf(&sb, "%d\n", q)
		}

		expected := solveQueries(arr, queries)
		got, err := runCase(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		tokens := strings.Fields(got)
		if len(tokens) != len(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\n", caseNum, len(expected), len(tokens))
			os.Exit(1)
		}
		for i, val := range expected {
			if tokens[i] != strconv.Itoa(val) {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", caseNum, val, tokens[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
