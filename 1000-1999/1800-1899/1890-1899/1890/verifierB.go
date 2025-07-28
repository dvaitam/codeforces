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

func isGood(s string) bool {
	for i := 0; i+1 < len(s); i++ {
		if s[i] == s[i+1] {
			return false
		}
	}
	return true
}

func expectedB(s, t string) string {
	if isGood(s) {
		return "YES"
	}
	if !isGood(t) {
		return "NO"
	}
	has00 := strings.Contains(s, "00")
	has11 := strings.Contains(s, "11")
	if has00 && has11 {
		return "NO"
	}
	first := t[0]
	last := t[len(t)-1]
	if has00 {
		if first == '1' && last == '1' {
			return "YES"
		}
		return "NO"
	}
	if has11 {
		if first == '0' && last == '0' {
			return "YES"
		}
		return "NO"
	}
	return "NO"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not read testcasesB.txt:", err)
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
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
			os.Exit(1)
		}
		s := scan.Text()
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
			os.Exit(1)
		}
		tStr := scan.Text()
		var input strings.Builder
		fmt.Fprintf(&input, "1\n%d %d\n%s\n%s\n", n, m, s, tStr)
		expect := expectedB(s, tStr)
		var cmd *exec.Cmd
		if strings.HasSuffix(bin, ".go") {
			cmd = exec.Command("go", "run", bin)
		} else {
			cmd = exec.Command(bin)
		}
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", caseNum, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", caseNum, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
