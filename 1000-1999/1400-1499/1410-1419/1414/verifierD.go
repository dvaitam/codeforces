package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func isPalindrome(s string) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanLines)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscanf(scanner.Text(), "%d", &t)
	inputs := make([]string, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Printf("missing input for test %d\n", i+1)
			os.Exit(1)
		}
		s := scanner.Text()
		inputs[i] = s
		if isPalindrome(s) {
			expected[i] = "Yes"
		} else {
			expected[i] = "No"
		}
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("execution failed: %v\n", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
