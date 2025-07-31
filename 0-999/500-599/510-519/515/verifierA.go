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

func canReach(a, b, s int64) bool {
	d := abs(a) + abs(b)
	return s >= d && (s-d)%2 == 0
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
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
	tests := make([][3]int64, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		a, _ := strconv.ParseInt(scan.Text(), 10, 64)
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		b, _ := strconv.ParseInt(scan.Text(), 10, 64)
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		s, _ := strconv.ParseInt(scan.Text(), 10, 64)
		tests[i] = [3]int64{a, b, s}
		if canReach(a, b, s) {
			expected[i] = "Yes"
		} else {
			expected[i] = "No"
		}
	}
	for i := 0; i < t; i++ {
		in := fmt.Sprintf("%d %d %d\n", tests[i][0], tests[i][1], tests[i][2])
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = strings.NewReader(in)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("execution failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		outScan := bufio.NewScanner(bytes.NewReader(out))
		outScan.Split(bufio.ScanWords)
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		if outScan.Text() != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], outScan.Text())
			os.Exit(1)
		}
		if outScan.Scan() {
			fmt.Printf("extra output detected in test %d\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
