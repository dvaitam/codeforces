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

func expected(n, m int) string {
	var sb strings.Builder
	for i := 2; i <= n; i += 2 {
		base := (i - 1) * m
		for j := 1; j <= m; j++ {
			sb.WriteString(strconv.Itoa(base + j))
			if j < m {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	for i := 1; i <= n; i += 2 {
		base := (i - 1) * m
		for j := 1; j <= m; j++ {
			sb.WriteString(strconv.Itoa(base + j))
			if j < m {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not read testcasesC.txt:", err)
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
		input := fmt.Sprintf("1\n%d %d\n", n, m)
		exp := expected(n, m)
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
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", caseNum, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\n", caseNum, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
