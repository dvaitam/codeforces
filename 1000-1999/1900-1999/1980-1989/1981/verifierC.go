package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

const maxVal = 5

func valid(a, b []int) bool {
	n := len(a)
	for i := 0; i < n; i++ {
		if a[i] != -1 && a[i] != b[i] {
			return false
		}
	}
	for i := 0; i < n-1; i++ {
		if !(b[i] == b[i+1]/2 || b[i+1] == b[i]/2) {
			return false
		}
	}
	return true
}

func findSolution(a []int) ([]int, bool) {
	n := len(a)
	b := make([]int, n)
	var res []int
	var dfs func(int)
	dfs = func(pos int) {
		if res != nil {
			return
		}
		if pos == n {
			if valid(a, b) {
				res = append([]int(nil), b...)
			}
			return
		}
		if a[pos] != -1 {
			b[pos] = a[pos]
			dfs(pos + 1)
			return
		}
		for v := 1; v <= maxVal && res == nil; v++ {
			b[pos] = v
			dfs(pos + 1)
		}
	}
	dfs(0)
	if res == nil {
		return nil, false
	}
	return res, true
}

func parseOutput(scan *bufio.Scanner, n int) ([]int, bool) {
	var b []int
	for i := 0; i < n; i++ {
		if !scan.Scan() {
			return nil, false
		}
		v, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, false
		}
		b = append(b, v)
	}
	return b, true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([][]int, t)
	expectSolution := make([]bool, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		a := make([]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			a[j] = v
		}
		if _, ok := findSolution(a); ok {
			expectSolution[i] = true
		} else {
			expectSolution[i] = false
		}
		cases[i] = a
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		a := cases[i]
		n := len(a)
		if expectSolution[i] {
			b, ok := parseOutput(outScan, n)
			if !ok {
				fmt.Printf("bad output for test %d\n", i+1)
				os.Exit(1)
			}
			if !valid(a, b) {
				fmt.Printf("test %d failed: sequence invalid\n", i+1)
				os.Exit(1)
			}
		} else {
			if !outScan.Scan() {
				fmt.Printf("missing output for test %d\n", i+1)
				os.Exit(1)
			}
			if outScan.Text() != "-1" {
				fmt.Printf("test %d failed: expected -1\n", i+1)
				os.Exit(1)
			}
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
