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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	// parse test cases
	in := bufio.NewScanner(bytes.NewReader(data))
	if !in.Scan() {
		fmt.Println("empty test file")
		os.Exit(1)
	}
	t, err := strconv.Atoi(strings.TrimSpace(in.Text()))
	if err != nil {
		fmt.Println("invalid test count:", err)
		os.Exit(1)
	}
	cases := make([]int, 0, t)
	for in.Scan() {
		n, err := strconv.Atoi(strings.TrimSpace(in.Text()))
		if err != nil {
			fmt.Println("invalid n value:", err)
			os.Exit(1)
		}
		cases = append(cases, n)
	}
	if len(cases) != t {
		fmt.Println("number of tests does not match")
		os.Exit(1)
	}

	bin := os.Args[1]
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("execution failed: %v\n%s", err, string(out))
		os.Exit(1)
	}

	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for idx, n := range cases {
		sum := 0
		for i := 1; i <= n; i++ {
			if !outScan.Scan() {
				fmt.Printf("missing output for test %d\n", idx+1)
				os.Exit(1)
			}
			v, err := strconv.Atoi(outScan.Text())
			if err != nil {
				fmt.Printf("invalid integer for test %d: %v\n", idx+1, err)
				os.Exit(1)
			}
			if v < 1 || v > 1000 {
				fmt.Printf("value out of range in test %d: %d\n", idx+1, v)
				os.Exit(1)
			}
			if v%i != 0 {
				fmt.Printf("value %d not divisible by index %d in test %d\n", v, i, idx+1)
				os.Exit(1)
			}
			sum += v
		}
		if sum%n != 0 {
			fmt.Printf("sum not divisible by n in test %d\n", idx+1)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
