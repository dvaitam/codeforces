package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solveCase(n, m, k int, heights [][]int64, types []string) string {
	ones := make([][]int, n)
	var sum0, sum1 int64
	for i := 0; i < n; i++ {
		ones[i] = make([]int, m)
		for j := 0; j < m; j++ {
			if types[i][j] == '1' {
				ones[i][j] = 1
				sum1 += heights[i][j]
			} else {
				sum0 += heights[i][j]
			}
		}
	}
	pref := make([][]int, n+1)
	for i := range pref {
		pref[i] = make([]int, m+1)
	}
	for i := 0; i < n; i++ {
		row := 0
		for j := 0; j < m; j++ {
			row += ones[i][j]
			pref[i+1][j+1] = pref[i][j+1] + row
		}
	}
	var g int64
	for i := 0; i+k <= n; i++ {
		for j := 0; j+k <= m; j++ {
			onesCount := pref[i+k][j+k] - pref[i][j+k] - pref[i+k][j] + pref[i][j]
			diff := int64(k*k - 2*onesCount)
			if diff < 0 {
				diff = -diff
			}
			g = gcd(g, diff)
		}
	}
	diffTotal := sum1 - sum0
	if g == 0 {
		if diffTotal == 0 {
			return "YES"
		}
		return "NO"
	}
	if diffTotal%g == 0 {
		return "YES"
	}
	return "NO"
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
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]string, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		nVal, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		mVal, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		kVal, _ := strconv.Atoi(scan.Text())
		heights := make([][]int64, nVal)
		for i := 0; i < nVal; i++ {
			heights[i] = make([]int64, mVal)
			for j := 0; j < mVal; j++ {
				scan.Scan()
				v, _ := strconv.ParseInt(scan.Text(), 10, 64)
				heights[i][j] = v
			}
		}
		types := make([]string, nVal)
		for i := 0; i < nVal; i++ {
			scan.Scan()
			types[i] = scan.Text()
		}
		expected[caseIdx] = solveCase(nVal, mVal, kVal, heights, types)
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
