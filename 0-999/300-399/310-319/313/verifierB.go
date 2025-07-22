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

func solveCase(s string, queries [][2]int) []int {
	n := len(s)
	p := make([]int, n+1)
	for i := 1; i < n; i++ {
		p[i+1] = p[i]
		if s[i] == s[i-1] {
			p[i+1]++
		}
	}
	res := make([]int, len(queries))
	for i, q := range queries {
		l, r := q[0], q[1]
		res[i] = p[r] - p[l]
	}
	return res
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read testcases: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	scan.Scan()
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		s := scan.Text()
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		m, _ := strconv.Atoi(scan.Text())
		queries := make([][2]int, m)
		for i := 0; i < m; i++ {
			scan.Scan()
			l, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			r, _ := strconv.Atoi(scan.Text())
			queries[i] = [2]int{l, r}
		}
		expectedSlice := solveCase(s, queries)
		var input strings.Builder
		input.WriteString(s)
		input.WriteByte('\n')
		input.WriteString(strconv.Itoa(m))
		input.WriteByte('\n')
		for _, q := range queries {
			input.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
		}
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		expectedStrs := make([]string, len(expectedSlice))
		for i, v := range expectedSlice {
			expectedStrs[i] = strconv.Itoa(v)
		}
		expected := strings.Join(expectedStrs, "\n")
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\n\ngot:\n%s\n", caseNum, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
