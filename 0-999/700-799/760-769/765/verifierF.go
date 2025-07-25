package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func solveCase(n int, arr []int, queries [][2]int) []int {
	res := make([]int, len(queries))
	for idx, q := range queries {
		l, r := q[0], q[1]
		b := make([]int, r-l+1)
		copy(b, arr[l-1:r])
		sort.Ints(b)
		minDiff := b[1] - b[0]
		for i := 2; i < len(b); i++ {
			d := b[i] - b[i-1]
			if d < minDiff {
				minDiff = d
			}
		}
		res[idx] = minDiff
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
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(scan.Text(), &t)
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		var n int
		fmt.Sscan(scan.Text(), &n)
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		fields := strings.Fields(scan.Text())
		if len(fields) != n {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		arr := make([]int, n)
		for i, v := range fields {
			fmt.Sscan(v, &arr[i])
		}
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		var m int
		fmt.Sscan(scan.Text(), &m)
		queries := make([][2]int, m)
		for i := 0; i < m; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &queries[i][0], &queries[i][1])
		}
		expected := solveCase(n, arr, queries)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", arr[i])
		}
		input.WriteByte('\n')
		fmt.Fprintf(&input, "%d\n", m)
		for i := 0; i < m; i++ {
			fmt.Fprintf(&input, "%d %d\n", queries[i][0], queries[i][1])
		}
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", caseNum, err)
			os.Exit(1)
		}
		outLines := strings.Fields(strings.TrimSpace(string(out)))
		if len(outLines) != len(expected) {
			fmt.Printf("case %d: expected %d lines got %d\n", caseNum, len(expected), len(outLines))
			os.Exit(1)
		}
		for i, exp := range expected {
			var val int
			fmt.Sscan(outLines[i], &val)
			if val != exp {
				fmt.Printf("case %d line %d: expected %d got %d\n", caseNum, i+1, exp, val)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed!")
}
