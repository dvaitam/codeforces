package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func solve(costs []int64) int64 {
	n := len(costs)
	minOdd := costs[0]
	minEven := costs[1]
	sumOdd := costs[0]
	sumEven := costs[1]
	cntOdd := 1
	cntEven := 1
	ans := (int64(n)-int64(cntOdd))*minOdd + (int64(n)-int64(cntEven))*minEven + sumOdd + sumEven
	for i := 2; i < n; i++ {
		v := costs[i]
		if i%2 == 0 {
			sumOdd += v
			cntOdd++
			if v < minOdd {
				minOdd = v
			}
		} else {
			sumEven += v
			cntEven++
			if v < minEven {
				minEven = v
			}
		}
		cur := (int64(n)-int64(cntOdd))*minOdd + (int64(n)-int64(cntEven))*minEven + sumOdd + sumEven
		if cur < ans {
			ans = cur
		}
	}
	return ans
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
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		costs := make([]int64, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			costs[j] = v
		}
		expected[i] = fmt.Sprintf("%d", solve(costs))
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
