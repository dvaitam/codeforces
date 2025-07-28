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

const INF int64 = 1e18

func solveCase(n int, pos []int, val []int64) []int64 {
	ans := make([]int64, n)
	for i := 0; i < n; i++ {
		ans[i] = INF
	}
	for i := 0; i < len(pos); i++ {
		idx := pos[i] - 1
		if val[i] < ans[idx] {
			ans[idx] = val[i]
		}
	}
	for i := 1; i < n; i++ {
		if ans[i-1]+1 < ans[i] {
			ans[i] = ans[i-1] + 1
		}
	}
	for i := n - 2; i >= 0; i-- {
		if ans[i+1]+1 < ans[i] {
			ans[i] = ans[i+1] + 1
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		pos := make([]int, k)
		for j := 0; j < k; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			pos[j], _ = strconv.Atoi(scan.Text())
		}
		val := make([]int64, k)
		for j := 0; j < k; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			v, _ := strconv.Atoi(scan.Text())
			val[j] = int64(v)
		}
		ans := solveCase(n, pos, val)
		strs := make([]string, n)
		for j, v := range ans {
			strs[j] = strconv.FormatInt(v, 10)
		}
		expected[i] = strings.Join(strs, " ")
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
