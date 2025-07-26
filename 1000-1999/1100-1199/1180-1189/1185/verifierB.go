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

func match(s, t string) bool {
	i, j := 0, 0
	for i < len(s) && j < len(t) {
		if s[i] != t[j] {
			return false
		}
		ch := s[i]
		cs, ct := 0, 0
		for i+cs < len(s) && s[i+cs] == ch {
			cs++
		}
		for j+ct < len(t) && t[j+ct] == ch {
			ct++
		}
		if ct < cs {
			return false
		}
		i += cs
		j += ct
	}
	return i == len(s) && j == len(t)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not open testcasesB.txt:", err)
		os.Exit(1)
	}
	// compute expected results
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	n, _ := strconv.Atoi(scanner.Text())
	expected := make([]string, n)
	// parse again line by line for expected
	scanner = bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan() // skip n
	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		s := scanner.Text()
		scanner.Scan()
		t := scanner.Text()
		if match(s, t) {
			expected[i] = "YES"
		} else {
			expected[i] = "NO"
		}
	}
	// run candidate
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outs := strings.Fields(string(out))
	if len(outs) != n {
		fmt.Printf("expected %d lines of output, got %d\n", n, len(outs))
		os.Exit(1)
	}
	for i := 0; i < n; i++ {
		if outs[i] != expected[i] {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, expected[i], outs[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
