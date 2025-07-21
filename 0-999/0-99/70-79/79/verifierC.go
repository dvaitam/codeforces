package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func ok(s string, b []string, l []int, k int) int {
	mx := -1
	for t := 0; t < len(b); t++ {
		start := k - l[t] + 1
		if start < 0 {
			continue
		}
		match := true
		for i := 0; i < l[t]; i++ {
			if s[start+i] != b[t][i] {
				match = false
				break
			}
		}
		if match && start > mx {
			mx = start
		}
	}
	return mx
}

func solveCaseC(s string, b []string) (int, int) {
	l := make([]int, len(b))
	for i := range b {
		l[i] = len(b[i])
	}
	p, mx, st := 0, -1, 0
	for i := 0; i < len(s); i++ {
		k := ok(s, b, l, i)
		if k != -1 {
			if i-p > mx {
				mx = i - p
				st = p
			}
			if k+1 > p {
				p = k + 1
			}
		}
	}
	if len(s)-p > mx {
		mx = len(s) - p
		st = p
	}
	return mx, st
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	T, _ := strconv.Atoi(scan.Text())
	cases := make([]struct {
		s string
		b []string
	}, T)
	expected := make([][2]int, T)
	for tc := 0; tc < T; tc++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		s := scan.Text()
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		b := make([]string, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			b[i] = scan.Text()
		}
		cases[tc] = struct {
			s string
			b []string
		}{s: s, b: b}
		mx, st := solveCaseC(s, b)
		expected[tc] = [2]int{mx, st}
	}
	for i, c := range cases {
		var buf bytes.Buffer
		fmt.Fprintln(&buf, c.s)
		fmt.Fprintln(&buf, len(c.b))
		for _, x := range c.b {
			fmt.Fprintln(&buf, x)
		}
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("execution failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		outScan := bufio.NewScanner(bytes.NewReader(out))
		outScan.Split(bufio.ScanWords)
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		mx, _ := strconv.Atoi(outScan.Text())
		if !outScan.Scan() {
			fmt.Printf("missing second value for test %d\n", i+1)
			os.Exit(1)
		}
		st, _ := strconv.Atoi(outScan.Text())
		if mx != expected[i][0] || st != expected[i][1] {
			fmt.Printf("test %d failed: expected %d %d got %d %d\n", i+1, expected[i][0], expected[i][1], mx, st)
			os.Exit(1)
		}
		if outScan.Scan() {
			fmt.Printf("extra output detected on case %d\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
