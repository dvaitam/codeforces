package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func solve(s string) string {
	n := len(s)
	dp := make([][][][]bool, n+1)
	for i := range dp {
		dp[i] = make([][][]bool, 2)
		for p := range dp[i] {
			dp[i][p] = make([][]bool, 2)
			for r := range dp[i][p] {
				dp[i][p][r] = make([]bool, 2)
			}
		}
	}
	dp[0][0][0][0] = true
	for i := 0; i < n; i++ {
		ch := s[i]
		for phase := 0; phase < 2; phase++ {
			for prev := 0; prev < 2; prev++ {
				for rem := 0; rem < 2; rem++ {
					if !dp[i][phase][prev][rem] {
						continue
					}
					if phase == 0 {
						if ch == '0' {
							dp[i+1][0][0][rem] = true
						} else {
							dp[i+1][1][0][rem] = true
						}
					} else {
						if ch == '1' {
							dp[i+1][1][0][rem] = true
						}
					}
					if prev == 0 {
						dp[i+1][phase][1][rem|1] = true
					}
				}
			}
		}
	}
	for phase := 0; phase < 2; phase++ {
		for prev := 0; prev < 2; prev++ {
			if dp[n][phase][prev][1] {
				return "YES"
			}
		}
	}
	return "NO"
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
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		s := scan.Text()
		expected[i] = solve(s)
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
