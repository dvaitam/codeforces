package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func max4(a, b, c, d int) int {
	if b > a {
		a = b
	}
	if c > a {
		a = c
	}
	if d > a {
		a = d
	}
	return a
}

func min4(a, b, c, d int) int {
	if b < a {
		a = b
	}
	if c < a {
		a = c
	}
	if d < a {
		a = d
	}
	return a
}

func expected(n, a, b, c, d int) int64 {
	low := max4(1, 1-(b-c), 1-(a-d), 1-(a+b-c-d))
	high := min4(n, n-(b-c), n-(a-d), n-(a+b-c-d))
	if low > high {
		return 0
	}
	return int64(high-low+1) * int64(n)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
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
	expect := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		aVal, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		bVal, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		cVal, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		dVal, _ := strconv.Atoi(scan.Text())
		ans := expected(n, aVal, bVal, cVal, dVal)
		expect[i] = fmt.Sprintf("%d", ans)
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
		if got != expect[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expect[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
