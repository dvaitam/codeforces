package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
)

func solveCase(nums []int) float64 {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	return float64(sum) / float64(len(nums))
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
		fmt.Println("invalid file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]float64, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		nums := make([]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			nums[j], _ = strconv.Atoi(scan.Text())
		}
		expected[i] = solveCase(nums)
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
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		v, err := strconv.ParseFloat(outScan.Text(), 64)
		if err != nil {
			fmt.Printf("bad output for test %d\n", i+1)
			os.Exit(1)
		}
		if math.Abs(v-expected[i]) > 1e-4 {
			fmt.Printf("test %d failed: expected %.6f got %.6f\n", i+1, expected[i], v)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
