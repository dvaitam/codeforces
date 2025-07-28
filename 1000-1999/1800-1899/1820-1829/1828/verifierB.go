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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	if !scan.Scan() {
		fmt.Println("empty test file")
		os.Exit(1)
	}
	t, err := strconv.Atoi(strings.TrimSpace(scan.Text()))
	if err != nil {
		fmt.Println("invalid test count:", err)
		os.Exit(1)
	}
	type test struct {
		n   int
		arr []int
	}
	tests := make([]test, 0, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("missing n value")
			os.Exit(1)
		}
		n, err := strconv.Atoi(strings.TrimSpace(scan.Text()))
		if err != nil {
			fmt.Println("invalid n:", err)
			os.Exit(1)
		}
		if !scan.Scan() {
			fmt.Println("missing permutation")
			os.Exit(1)
		}
		parts := strings.Fields(scan.Text())
		if len(parts) != n {
			fmt.Printf("invalid permutation length for n=%d\n", n)
			os.Exit(1)
		}
		arr := make([]int, n)
		for j, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil {
				fmt.Println("invalid integer in permutation:", err)
				os.Exit(1)
			}
			arr[j] = v
		}
		tests = append(tests, test{n: n, arr: arr})
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
	for idx, tc := range tests {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", idx+1)
			os.Exit(1)
		}
		got, err := strconv.Atoi(outScan.Text())
		if err != nil {
			fmt.Printf("invalid integer for test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		g := 0
		for i, v := range tc.arr {
			diff := v - (i + 1)
			if diff < 0 {
				diff = -diff
			}
			g = gcd(g, diff)
		}
		if got != g {
			fmt.Printf("test %d failed: expected %d got %d\n", idx+1, g, got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
