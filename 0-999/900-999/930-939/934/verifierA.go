package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
)

func solveCase(n, m int, a, b []int64) int64 {
	best := make([]int64, n)
	for i := 0; i < n; i++ {
		maxVal := int64(math.MinInt64)
		for j := 0; j < m; j++ {
			prod := a[i] * b[j]
			if prod > maxVal {
				maxVal = prod
			}
		}
		best[i] = maxVal
	}
	ans := int64(math.MaxInt64)
	for hide := 0; hide < n; hide++ {
		maxVal := int64(math.MinInt64)
		for i := 0; i < n; i++ {
			if i == hide {
				continue
			}
			if best[i] > maxVal {
				maxVal = best[i]
			}
		}
		if maxVal < ans {
			ans = maxVal
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			val, _ := strconv.ParseInt(scan.Text(), 10, 64)
			a[i] = val
		}
		b := make([]int64, m)
		for j := 0; j < m; j++ {
			scan.Scan()
			val, _ := strconv.ParseInt(scan.Text(), 10, 64)
			b[j] = val
		}
		expected := solveCase(n, m, a, b)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", a[i])
		}
		input.WriteByte('\n')
		for j := 0; j < m; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", b[j])
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n%s", caseNum, err, out.String())
			os.Exit(1)
		}
		outScan := bufio.NewScanner(bytes.NewReader(out.Bytes()))
		outScan.Split(bufio.ScanWords)
		if !outScan.Scan() {
			fmt.Printf("case %d: no output\n", caseNum)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(outScan.Text(), 10, 64)
		if err != nil {
			fmt.Printf("case %d: invalid output: %v\n", caseNum, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed: expected %d got %d\n", caseNum, expected, got)
			os.Exit(1)
		}
		if outScan.Scan() {
			fmt.Printf("case %d: extra output detected\n", caseNum)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
