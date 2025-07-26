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
	if n <= 2 {
		return "No"
	}
	var sb strings.Builder
	sb.WriteString("Yes\n")
	if n%2 == 0 {
		sb.WriteString(fmt.Sprintf("2 1 %d\n", n))
		sb.WriteString(fmt.Sprintf("%d", n-2))
		for i := 2; i < n; i++ {
			sb.WriteString(fmt.Sprintf(" %d", i))
		}
		sb.WriteString("\n")
	} else {
		k := (n + 1) / 2
		sb.WriteString(fmt.Sprintf("1 %d\n", k))
		sb.WriteString(fmt.Sprintf("%d", n-1))
		for i := 1; i <= n; i++ {
			if i == k {
				continue
			}
			sb.WriteString(fmt.Sprintf(" %d", i))
		}
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		input := fmt.Sprintf("%d\n", n)
		exp := expected(n)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
