package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func cwDist(pref []int64, total int64, i, j int) int64 {
	if i <= j {
		return pref[j-1] - pref[i-1]
	}
	return total - (pref[i-1] - pref[j-1])
}

func solveQuery(n int, d, h []int64, a, b int) int64 {
	banned := make(map[int]bool)
	for x := a; ; x = x%n + 1 {
		banned[x] = true
		if x == b {
			break
		}
	}
	open := []int{}
	for x := b%n + 1; x != a; x = x%n + 1 {
		if !banned[x] {
			open = append(open, x)
		}
	}
	// Build prefix distances along circle
	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + d[i-1]
	}
	total := pref[n]
	// Build prefix along open arc
	arcPref := make([]int64, len(open))
	for i := 1; i < len(open); i++ {
		arcPref[i] = arcPref[i-1] + cwDist(pref, total, open[i-1], open[i])
	}
	best := int64(-1 << 60)
	for i := 0; i < len(open); i++ {
		for j := i + 1; j < len(open); j++ {
			dist := arcPref[j] - arcPref[i]
			energy := dist + 2*(h[open[i]-1]+h[open[j]-1])
			if energy > best {
				best = energy
			}
		}
	}
	return best
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
	answers := make([][]int64, t)
	casesN := make([]int, t)
	for idx := 0; idx < t; idx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		d := make([]int64, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			val, _ := strconv.ParseInt(scan.Text(), 10, 64)
			d[i] = val
		}
		h := make([]int64, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			val, _ := strconv.ParseInt(scan.Text(), 10, 64)
			h[i] = val
		}
		res := make([]int64, m)
		for q := 0; q < m; q++ {
			scan.Scan()
			a, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			b, _ := strconv.Atoi(scan.Text())
			res[q] = solveQuery(n, d, h, a, b)
		}
		answers[idx] = res
		casesN[idx] = m
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
	for idx := 0; idx < t; idx++ {
		for j := 0; j < casesN[idx]; j++ {
			if !outScan.Scan() {
				fmt.Printf("missing output for test %d query %d\n", idx+1, j+1)
				os.Exit(1)
			}
			got, err := strconv.ParseInt(outScan.Text(), 10, 64)
			if err != nil {
				fmt.Printf("bad output for test %d query %d\n", idx+1, j+1)
				os.Exit(1)
			}
			if got != answers[idx][j] {
				fmt.Printf("test %d query %d failed: expected %d got %d\n", idx+1, j+1, answers[idx][j], got)
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
