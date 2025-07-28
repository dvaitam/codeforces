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

func solveCase(k int, a, b []int) ([]int, bool) {
	res := make([]int, 0, len(a)+len(b))
	i, j := 0, 0
	lines := k
	for i < len(a) || j < len(b) {
		if i < len(a) && a[i] == 0 {
			res = append(res, 0)
			lines++
			i++
			continue
		}
		if j < len(b) && b[j] == 0 {
			res = append(res, 0)
			lines++
			j++
			continue
		}
		if i < len(a) && a[i] <= lines {
			res = append(res, a[i])
			i++
			continue
		}
		if j < len(b) && b[j] <= lines {
			res = append(res, b[j])
			j++
			continue
		}
		return nil, false
	}
	return res, true
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
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		k, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		a := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			a[j], _ = strconv.Atoi(scan.Text())
		}
		b := make([]int, m)
		for j := 0; j < m; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			b[j], _ = strconv.Atoi(scan.Text())
		}
		if ans, ok := solveCase(k, a, b); ok {
			strs := make([]string, len(ans))
			for idx, v := range ans {
				strs[idx] = strconv.Itoa(v)
			}
			expected[i] = strings.Join(strs, " ")
		} else {
			expected[i] = "-1"
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
	outScan.Split(bufio.ScanLines)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := strings.TrimSpace(outScan.Text())
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
