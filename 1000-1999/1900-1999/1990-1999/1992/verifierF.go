package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func solve(n, x int, arr []int) int {
	ans := 1
	s := map[int]struct{}{1: {}}
	for _, v := range arr {
		s2 := make(map[int]struct{})
		for val := range s {
			if x%(val*v) == 0 {
				s2[val*v] = struct{}{}
			}
		}
		for k := range s2 {
			s[k] = struct{}{}
		}
		if _, ok := s[x]; ok {
			ans++
			s = map[int]struct{}{1: {}, v: {}}
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("problemF.txt")
	if err != nil {
		fmt.Println("could not read problemF.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]int, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		x, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			arr[j] = v
		}
		expected[i] = solve(n, x, arr)
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
		got, _ := strconv.Atoi(outScan.Text())
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %d got %d\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
