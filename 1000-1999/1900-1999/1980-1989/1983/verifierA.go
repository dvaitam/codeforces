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

func expected(n int) string {
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(i))
	}
	return sb.String()
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("failed to open testcasesA.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid testcases file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for idx := 1; idx <= t; idx++ {
		if !scan.Scan() {
			fmt.Printf("missing n for case %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		expect := expected(n)
		input := fmt.Sprintf("1\n%d\n", n)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %q got %q\n", idx, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
