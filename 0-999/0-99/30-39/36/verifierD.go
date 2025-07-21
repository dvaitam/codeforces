package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func win(x, y, k int, memo map[[2]int]bool) bool {
	if x == 0 && y == 0 {
		return false
	}
	key := [2]int{x, y}
	if v, ok := memo[key]; ok {
		return v
	}
	moves := []bool{}
	if x > 0 {
		moves = append(moves, win(x-1, y, k, memo))
	}
	if y > 0 {
		moves = append(moves, win(x, y-1, k, memo))
	}
	if x >= k && y >= k {
		moves = append(moves, win(x-k, y-k, k, memo))
	}
	res := !allTrue(moves)
	memo[key] = res
	return res
}

func allTrue(b []bool) bool {
	for _, v := range b {
		if !v {
			return false
		}
	}
	return true
}

func solveCase(n, m, k int) string {
	memo := make(map[[2]int]bool)
	if win(n-1, m-1, k, memo) {
		return "+"
	}
	return "-"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		fmt.Println("bad file")
		os.Exit(1)
	}
	var t, k int
	fmt.Sscan(scan.Text(), &t)
	if !scan.Scan() {
		fmt.Println("bad file")
		os.Exit(1)
	}
	fmt.Sscan(scan.Text(), &k)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		var n int
		fmt.Sscan(scan.Text(), &n)
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		var m int
		fmt.Sscan(scan.Text(), &m)
		expected[i] = solveCase(n, m, k)
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
			fmt.Printf("missing output for case %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expected[i] {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
