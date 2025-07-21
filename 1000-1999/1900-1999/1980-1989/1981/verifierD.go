package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func minDistinct(n int) int {
	k := 1
	for {
		if k%2 == 1 {
			if k*(k+1)/2 >= n-1 {
				return k
			}
		} else {
			if k*(k+1)/2-k/2+1 >= n-1 {
				return k
			}
		}
		k++
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	ns := make([]int, t)
	ks := make([]int, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		ns[i] = n
		ks[i] = minDistinct(n)
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
		n := ns[i]
		k := ks[i]
		vals := make([]int, n)
		for j := 0; j < n; j++ {
			if !outScan.Scan() {
				fmt.Printf("missing output for test %d\n", i+1)
				os.Exit(1)
			}
			v, err := strconv.Atoi(outScan.Text())
			if err != nil {
				fmt.Printf("bad integer in test %d\n", i+1)
				os.Exit(1)
			}
			if v < 1 || v > 300000 {
				fmt.Printf("value out of range in test %d\n", i+1)
				os.Exit(1)
			}
			vals[j] = v
		}
		distinct := map[int]struct{}{}
		products := map[int]struct{}{}
		for j := 0; j < n-1; j++ {
			distinct[vals[j]] = struct{}{}
			distinct[vals[j+1]] = struct{}{}
			prod := vals[j] * vals[j+1]
			if _, ok := products[prod]; ok {
				fmt.Printf("duplicate product in test %d\n", i+1)
				os.Exit(1)
			}
			products[prod] = struct{}{}
		}
		if len(distinct) != k {
			fmt.Printf("distinct count mismatch in test %d: expected %d got %d\n", i+1, k, len(distinct))
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
