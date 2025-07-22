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

func expected(p string) int {
	n := len(p)
	if n == 0 {
		return 0
	}
	leaves := 1
	for i := 1; i < n; i++ {
		if p[i] != '0' {
			leaves++
		}
	}
	if n >= 2 && p[1] != '0' && p[0] < p[1] {
		leaves--
	}
	return leaves
}

func runCase(bin, p string) error {
	input := p + "\n"
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	exp := expected(p)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "bad test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scanner.Text())
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "test %d: missing input\n", i+1)
			os.Exit(1)
		}
		p := scanner.Text()
		if err := runCase(bin, p); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra data in test file")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}
