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

func expected(n, k int, s string) string {
	ops := 0
	for i := 0; i < n; {
		if s[i] == 'B' {
			ops++
			i += k
		} else {
			i++
		}
	}
	return fmt.Sprintf("%d\n", ops)
}

func runCase(exe, input, expect string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect = strings.TrimSpace(expect)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		parts := strings.Fields(scan.Text())
		if len(parts) != 2 {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		k, _ := strconv.Atoi(parts[1])
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		s := strings.TrimSpace(scan.Text())
		input := fmt.Sprintf("1\n%d %d\n%s\n", n, k, s)
		exp := expected(n, k, s)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
