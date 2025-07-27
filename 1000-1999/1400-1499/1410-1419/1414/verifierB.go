package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scanner.Text())
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Printf("missing a for test %d\n", i+1)
			os.Exit(1)
		}
		aVal, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		if !scanner.Scan() {
			fmt.Printf("missing b for test %d\n", i+1)
			os.Exit(1)
		}
		bVal, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		expected[i] = fmt.Sprintf("%d", gcd(aVal, bVal))
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
