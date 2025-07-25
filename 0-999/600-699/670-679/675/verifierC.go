package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func expected(arr []int64) int {
	freq := make(map[int64]int)
	prefix := int64(0)
	maxFreq := 0
	for _, v := range arr {
		prefix += v
		freq[prefix]++
		if freq[prefix] > maxFreq {
			maxFreq = freq[prefix]
		}
	}
	return len(arr) - maxFreq
}

func main() {
	if len(os.Args) != 2 {
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
	expect := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			val, _ := strconv.ParseInt(scan.Text(), 10, 64)
			arr[j] = val
		}
		ans := expected(arr)
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
