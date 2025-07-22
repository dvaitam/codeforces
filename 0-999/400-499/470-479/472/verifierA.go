package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func isComposite(x int) bool {
	if x < 4 {
		return false
	}
	for i := 2; i*i <= x; i++ {
		if x%i == 0 {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
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
	inputs := make([]int, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		inputs[i] = n
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewReader(bytes.NewReader(out))
	for i := 0; i < t; i++ {
		var x, y int
		if _, err := fmt.Fscan(outScan, &x, &y); err != nil {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		n := inputs[i]
		if x+y != n || !isComposite(x) || !isComposite(y) {
			fmt.Printf("test %d failed: n=%d x=%d y=%d\n", i+1, n, x, y)
			os.Exit(1)
		}
	}
	if _, err := fmt.Fscan(outScan, new(int)); err == nil {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
