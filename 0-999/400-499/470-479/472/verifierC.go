package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func solveCase(n int, first, last []string, perm []int) string {
	small := make([]string, n)
	large := make([]string, n)
	for i := 0; i < n; i++ {
		if first[i] < last[i] {
			small[i], large[i] = first[i], last[i]
		} else {
			small[i], large[i] = last[i], first[i]
		}
	}
	prev := ""
	for _, idx := range perm {
		idx--
		opt1 := small[idx]
		opt2 := large[idx]
		var chosen string
		ok1 := opt1 > prev
		ok2 := opt2 > prev
		if ok1 && ok2 {
			if opt1 < opt2 {
				chosen = opt1
			} else {
				chosen = opt2
			}
		} else if ok1 {
			chosen = opt1
		} else if ok2 {
			chosen = opt2
		} else {
			return "NO"
		}
		prev = chosen
	}
	return "YES"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
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
	// We'll build input for binary same as in file
	input := data
	// For computing expected results, parse individually
	scan = bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	scan.Scan() // t
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		first := make([]string, n)
		last := make([]string, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			first[i] = scan.Text()
			scan.Scan()
			last[i] = scan.Text()
		}
		perm := make([]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			perm[i], _ = strconv.Atoi(scan.Text())
		}
		expected[caseIdx] = solveCase(n, first, last, perm)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	outBytes, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(outBytes))
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
