package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func maxWins(a []int64, l, r int64, idx int) int {
	if idx == len(a) {
		return 0
	}
	best := -1
	for j := idx; j < len(a); j++ {
		sum := int64(0)
		for t := idx; t <= j; t++ {
			sum += a[t]
		}
		win := 0
		if sum >= l && sum <= r {
			win = 1
		}
		val := win + maxWins(a, l, r, j+1)
		if val > best {
			best = val
		}
	}
	return best
}

func solveCase(n int, l, r int64, arr []int64) int {
	return maxWins(arr, l, r, 0)
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
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		lVal, _ := strconv.ParseInt(scan.Text(), 10, 64)
		scan.Scan()
		rVal, _ := strconv.ParseInt(scan.Text(), 10, 64)
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			arr[j] = v
		}
		expected[i] = fmt.Sprintf("%d", solveCase(n, lVal, rVal, arr))
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
