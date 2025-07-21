package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

const infB = int64(1 << 60)

func solveCaseB(n int, quals []int, edges [][3]int) (int64, bool) {
	best := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		best[i] = infB
	}
	for _, e := range edges {
		b, c := e[1], e[2]
		if int64(c) < best[b] {
			best[b] = int64(c)
		}
	}
	maxq := -1
	root := -1
	cnt := 0
	for i := 1; i <= n; i++ {
		if quals[i] > maxq {
			maxq = quals[i]
			root = i
			cnt = 1
		} else if quals[i] == maxq {
			cnt++
		}
	}
	if cnt != 1 {
		return 0, false
	}
	var sum int64
	for i := 1; i <= n; i++ {
		if i == root {
			continue
		}
		if best[i] == infB {
			return 0, false
		}
		sum += best[i]
	}
	return sum, true
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
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		quals := make([]int, n+1)
		for j := 1; j <= n; j++ {
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			quals[j] = v
		}
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		edges := make([][3]int, m)
		for j := 0; j < m; j++ {
			scan.Scan()
			a, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			b, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			c, _ := strconv.Atoi(scan.Text())
			edges[j] = [3]int{a, b, c}
		}
		if res, ok := solveCaseB(n, quals, edges); ok {
			expected[i] = fmt.Sprintf("%d", res)
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
