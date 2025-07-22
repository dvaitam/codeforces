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

func expected(n, m int, a, b []int) int {
	i, j, matched := 0, 0, 0
	for i < n && j < m {
		if b[j] >= a[i] {
			matched++
			i++
			j++
		} else {
			j++
		}
	}
	return n - matched
}

func runCase(bin string, n, m int, a, b []int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for idx, v := range a {
		if idx > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
	}
	input.WriteByte('\n')
	for idx, v := range b {
		if idx > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
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
	exp := expected(n, m, a, b)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
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
			fmt.Fprintf(os.Stderr, "test %d: missing n\n", i+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scanner.Text())
		scanner.Scan()
		m, _ := strconv.Atoi(scanner.Text())
		a := make([]int, n)
		for j := 0; j < n; j++ {
			scanner.Scan()
			val, _ := strconv.Atoi(scanner.Text())
			a[j] = val
		}
		b := make([]int, m)
		for j := 0; j < m; j++ {
			scanner.Scan()
			val, _ := strconv.Atoi(scanner.Text())
			b[j] = val
		}
		if err := runCase(bin, n, m, a, b); err != nil {
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
