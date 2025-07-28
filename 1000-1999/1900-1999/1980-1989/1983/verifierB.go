package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(n, m int, a, b []string) string {
	for i := 0; i < n; i++ {
		sum := 0
		for j := 0; j < m; j++ {
			ai := int(a[i][j] - '0')
			bi := int(b[i][j] - '0')
			diff := (bi - ai) % 3
			if diff < 0 {
				diff += 3
			}
			sum += diff
		}
		if sum%3 != 0 {
			return "NO"
		}
	}
	for j := 0; j < m; j++ {
		sum := 0
		for i := 0; i < n; i++ {
			ai := int(a[i][j] - '0')
			bi := int(b[i][j] - '0')
			diff := (bi - ai) % 3
			if diff < 0 {
				diff += 3
			}
			sum += diff
		}
		if sum%3 != 0 {
			return "NO"
		}
	}
	return "YES"
}

func runCandidate(bin, input string) (string, error) {
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
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("failed to open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		fmt.Println("invalid test count")
		os.Exit(1)
	}
	for caseNum := 1; caseNum <= t; caseNum++ {
		var n, m int
		if _, err := fmt.Fscan(in, &n, &m); err != nil {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		a := make([]string, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(in, &a[i]); err != nil {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
		}
		b := make([]string, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(in, &b[i]); err != nil {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
		}
		expect := expected(n, m, a, b)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			input.WriteString(a[i])
			input.WriteByte('\n')
		}
		for i := 0; i < n; i++ {
			input.WriteString(b[i])
			input.WriteByte('\n')
		}
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %q got %q\n", caseNum, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
