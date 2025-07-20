package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type pair struct{ a, b int }

func prefixRepeat(s string, d int) int {
	var b strings.Builder
	for b.Len() < d {
		b.WriteString(s)
	}
	str := b.String()[:d]
	val, _ := strconv.Atoi(str)
	return val
}

func solve(n int) []pair {
	ns := strconv.Itoa(n)
	l := len(ns)
	res := make([]pair, 0)
	for a := 1; a <= 10000; a++ {
		for d := 1; d <= 7; d++ {
			b := l*a - d
			if b < 1 || b > a*n || b > 10000 {
				continue
			}
			if d > l*a {
				continue
			}
			pref := prefixRepeat(ns, d)
			if pref == a*n-b {
				res = append(res, pair{a, b})
			}
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([][]pair, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		expected[i] = solve(n)
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
			fmt.Printf("missing count for test %d\n", i+1)
			os.Exit(1)
		}
		cnt, _ := strconv.Atoi(outScan.Text())
		if cnt != len(expected[i]) {
			fmt.Printf("test %d failed: expected count %d got %d\n", i+1, len(expected[i]), cnt)
			os.Exit(1)
		}
		for j := 0; j < cnt; j++ {
			if !outScan.Scan() {
				fmt.Printf("missing output for test %d pair %d\n", i+1, j+1)
				os.Exit(1)
			}
			a, _ := strconv.Atoi(outScan.Text())
			outScan.Scan()
			b, _ := strconv.Atoi(outScan.Text())
			if j >= len(expected[i]) || a != expected[i][j].a || b != expected[i][j].b {
				fmt.Printf("test %d pair %d mismatch\n", i+1, j+1)
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
