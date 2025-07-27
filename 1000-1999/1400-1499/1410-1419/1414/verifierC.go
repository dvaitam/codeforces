package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func fib(n int) int64 {
	if n <= 1 {
		return int64(n)
	}
	a, b := int64(0), int64(1)
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
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
			fmt.Printf("missing input for test %d\n", i+1)
			os.Exit(1)
		}
		nVal, _ := strconv.Atoi(scanner.Text())
		expected[i] = fmt.Sprintf("%d", fib(nVal))
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
