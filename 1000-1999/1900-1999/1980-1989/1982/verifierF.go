package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func findSegment(a []int) (int, int) {
	n := len(a)
	l := 0
	for l+1 < n && a[l] <= a[l+1] {
		l++
	}
	if l == n-1 {
		return -1, -1
	}
	r := n - 1
	for r > 0 && a[r-1] <= a[r] {
		r--
	}
	minV, maxV := a[l], a[l]
	for i := l; i <= r; i++ {
		if a[i] < minV {
			minV = a[i]
		}
		if a[i] > maxV {
			maxV = a[i]
		}
	}
	for l > 0 && a[l-1] > minV {
		l--
	}
	for r < n-1 && a[r+1] < maxV {
		r++
	}
	return l + 1, r + 1
}

func solveCase(n int, arr []int, queries [][2]int) []string {
	res := make([]string, len(queries)+1)
	l, r := findSegment(arr)
	res[0] = fmt.Sprintf("%d %d", l, r)
	for i, q := range queries {
		arr[q[0]-1] = q[1]
		l, r = findSegment(arr)
		res[i+1] = fmt.Sprintf("%d %d", l, r)
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Println("could not read testcasesF.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([][]string, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		scan.Scan()
		nVal, _ := strconv.Atoi(scan.Text())
		arr := make([]int, nVal)
		for i := 0; i < nVal; i++ {
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			arr[i] = v
		}
		scan.Scan()
		qVal, _ := strconv.Atoi(scan.Text())
		queries := make([][2]int, qVal)
		for i := 0; i < qVal; i++ {
			scan.Scan()
			p, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			queries[i] = [2]int{p, v}
		}
		expected[caseIdx] = solveCase(nVal, append([]int(nil), arr...), queries)
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
		for j := 0; j < len(expected[i]); j++ {
			if !outScan.Scan() {
				fmt.Printf("missing output for test %d line %d\n", i+1, j+1)
				os.Exit(1)
			}
			gotL := outScan.Text()
			if !outScan.Scan() {
				fmt.Printf("missing output for test %d line %d second number\n", i+1, j+1)
				os.Exit(1)
			}
			gotR := outScan.Text()
			expParts := expected[i][j]
			// expected string already "l r"
			var expL, expR string
			fmt.Sscan(expParts, &expL, &expR)
			if gotL != expL || gotR != expR {
				fmt.Printf("test %d step %d failed: expected %s got %s %s\n", i+1, j, expParts, gotL, gotR)
				os.Exit(1)
			}
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
