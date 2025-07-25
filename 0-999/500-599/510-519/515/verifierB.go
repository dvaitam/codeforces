package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func canAllHappy(n, m int, boys, girls []int) bool {
	g := gcd(n, m)
	seen := make([]bool, g)
	for _, x := range boys {
		seen[x%g] = true
	}
	for _, y := range girls {
		seen[y%g] = true
	}
	for i := 0; i < g; i++ {
		if !seen[i] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
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
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		bcnt, _ := strconv.Atoi(scan.Text())
		boys := make([]int, bcnt)
		for j := 0; j < bcnt; j++ {
			scan.Scan()
			val, _ := strconv.Atoi(scan.Text())
			boys[j] = val
		}
		scan.Scan()
		gcnt, _ := strconv.Atoi(scan.Text())
		girls := make([]int, gcnt)
		for j := 0; j < gcnt; j++ {
			scan.Scan()
			val, _ := strconv.Atoi(scan.Text())
			girls[j] = val
		}
		if canAllHappy(n, m, boys, girls) {
			expected[i] = "Yes"
		} else {
			expected[i] = "No"
		}
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
		if outScan.Text() != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], outScan.Text())
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
