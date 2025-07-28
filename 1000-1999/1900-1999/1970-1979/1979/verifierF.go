package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read testcasesF.txt: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	sizes := make([]int, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Fprintln(os.Stderr, "bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		sizes[i] = n
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "execution failed: %v\n", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		n := sizes[i]
		seen := make([]bool, n+1)
		for j := 0; j < n; j++ {
			if !outScan.Scan() {
				fmt.Fprintf(os.Stderr, "missing output for test %d\n", i+1)
				os.Exit(1)
			}
			val, err := strconv.Atoi(outScan.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "bad output for test %d\n", i+1)
				os.Exit(1)
			}
			if val < 1 || val > n {
				fmt.Fprintf(os.Stderr, "value out of range on test %d: %d\n", i+1, val)
				os.Exit(1)
			}
			if seen[val] {
				fmt.Fprintf(os.Stderr, "duplicate value on test %d: %d\n", i+1, val)
				os.Exit(1)
			}
			seen[val] = true
		}
		for v := 1; v <= n; v++ {
			if !seen[v] {
				fmt.Fprintf(os.Stderr, "missing value %d on test %d\n", v, i+1)
				os.Exit(1)
			}
		}
	}
	if outScan.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
