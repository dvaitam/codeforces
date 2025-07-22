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

func expectedNumber(a string, k int) string {
	s := []byte(a)
	n := len(s)
	for i := 0; i < n && k > 0; i++ {
		maxDigit := s[i]
		pos := i
		limit := i + k
		if limit >= n {
			limit = n - 1
		}
		for j := i + 1; j <= limit; j++ {
			if s[j] > maxDigit {
				maxDigit = s[j]
				pos = j
			}
		}
		if pos != i {
			for j := pos; j > i; j-- {
				s[j], s[j-1] = s[j-1], s[j]
			}
			k -= pos - i
		}
	}
	return string(s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "bad test case on line %d\n", idx)
			os.Exit(1)
		}
		a := parts[0]
		k, _ := strconv.Atoi(parts[1])
		expect := expectedNumber(a, k)
		input := fmt.Sprintf("%s %d\n", a, k)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
