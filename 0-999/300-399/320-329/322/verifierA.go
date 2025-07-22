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

func checkSchedule(m, n int, out string) error {
	scan := bufio.NewScanner(strings.NewReader(out))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return fmt.Errorf("missing pair count")
	}
	k, err := strconv.Atoi(scan.Text())
	if err != nil {
		return fmt.Errorf("invalid pair count: %v", err)
	}
	if k != m+n-1 {
		return fmt.Errorf("expected %d pairs got %d", m+n-1, k)
	}
	usedBoy := make([]bool, m+1)
	usedGirl := make([]bool, n+1)
	for i := 0; i < k; i++ {
		if !scan.Scan() {
			return fmt.Errorf("missing boy index for pair %d", i+1)
		}
		boy, err := strconv.Atoi(scan.Text())
		if err != nil {
			return fmt.Errorf("bad boy index on pair %d", i+1)
		}
		if !scan.Scan() {
			return fmt.Errorf("missing girl index for pair %d", i+1)
		}
		girl, err := strconv.Atoi(scan.Text())
		if err != nil {
			return fmt.Errorf("bad girl index on pair %d", i+1)
		}
		if boy < 1 || boy > m || girl < 1 || girl > n {
			return fmt.Errorf("pair %d out of range", i+1)
		}
		if usedBoy[boy] && usedGirl[girl] {
			return fmt.Errorf("pair %d violates rule", i+1)
		}
		usedBoy[boy] = true
		usedGirl[girl] = true
	}
	if scan.Scan() {
		return fmt.Errorf("extra output detected")
	}
	return nil
}

func runCase(bin string, m, n int) error {
	input := fmt.Sprintf("%d %d\n", m, n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return checkSchedule(m, n, out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "empty test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing m\n", i+1)
			os.Exit(1)
		}
		m, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing n\n", i+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if err := runCase(bin, m, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
