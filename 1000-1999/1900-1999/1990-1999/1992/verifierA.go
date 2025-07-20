package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func solve(a, b, c int) int {
	best := 0
	for i := 0; i <= 5; i++ {
		for j := 0; j <= 5-i; j++ {
			k := 5 - i - j
			prod := (a + i) * (b + j) * (c + k)
			if prod > best {
				best = prod
			}
		}
	}
	return best
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("problemA.txt")
	if err != nil {
		fmt.Println("could not read problemA.txt:", err)
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
			fmt.Println("not enough values in test file")
			os.Exit(1)
		}
		a, _ := strconv.Atoi(scanner.Text())
		scanner.Scan()
		b, _ := strconv.Atoi(scanner.Text())
		scanner.Scan()
		c, _ := strconv.Atoi(scanner.Text())
		expected[i] = fmt.Sprintf("%d", solve(a, b, c))
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
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
