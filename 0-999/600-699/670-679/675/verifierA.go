package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func expected(a, b, c int64) string {
	if a == b {
		return "YES"
	}
	if c == 0 {
		return "NO"
	}
	if (b-a)%c == 0 && (b-a)/c >= 0 {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
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
	expectedOut := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		aVal, _ := strconv.ParseInt(scan.Text(), 10, 64)
		scan.Scan()
		bVal, _ := strconv.ParseInt(scan.Text(), 10, 64)
		scan.Scan()
		cVal, _ := strconv.ParseInt(scan.Text(), 10, 64)
		expectedOut[i] = expected(aVal, bVal, cVal)
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("execution failed: %v\n%s", err, out)
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
		if got != expectedOut[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expectedOut[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
