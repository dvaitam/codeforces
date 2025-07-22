package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

func solveCase(n, k int, floors []int) int64 {
	sort.Ints(floors)
	var ans int64
	for i := n - 1; i >= 0; i -= k {
		ans += int64(2 * (floors[i] - 1))
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
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
	inputs := make([]struct {
		n, k   int
		floors []int
	}, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		floors := make([]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			floors[j] = v
		}
		inputs[i] = struct {
			n, k   int
			floors []int
		}{n, k, floors}
	}
	expected := make([]int64, t)
	for i, in := range inputs {
		expected[i] = solveCase(in.n, in.k, append([]int(nil), in.floors...))
	}
	cmd := exec.Command(bin)
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
		gotStr := outScan.Text()
		got, _ := strconv.ParseInt(gotStr, 10, 64)
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
