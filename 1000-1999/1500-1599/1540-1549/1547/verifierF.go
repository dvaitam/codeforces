package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func buildSparse(arr []int) ([][]int, []int) {
	n := len(arr)
	log := make([]int, n+1)
	for i := 2; i <= n; i++ {
		log[i] = log[i/2] + 1
	}
	K := log[n] + 1
	st := make([][]int, K)
	st[0] = make([]int, n)
	copy(st[0], arr)
	for k := 1; k < K; k++ {
		size := n - (1 << k) + 1
		st[k] = make([]int, size)
		step := 1 << (k - 1)
		for i := 0; i < size; i++ {
			st[k][i] = gcd(st[k-1][i], st[k-1][i+step])
		}
	}
	return st, log
}

func query(st [][]int, log []int, l, r int) int {
	length := r - l
	k := log[length]
	return gcd(st[k][l], st[k][r-(1<<k)])
}

func allEqual(arr []int) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i] != arr[0] {
			return false
		}
	}
	return true
}

func solveCase(a []int) int {
	n := len(a)
	if allEqual(a) {
		return 0
	}
	g := 0
	for _, v := range a {
		g = gcd(g, v)
	}
	b := make([]int, 2*n)
	copy(b, a)
	copy(b[n:], a)
	st, log := buildSparse(b)
	low, high := 1, n
	for low < high {
		mid := (low + high) / 2
		ok := true
		for i := 0; i < n; i++ {
			if query(st, log, i, i+mid) != g {
				ok = false
				break
			}
		}
		if ok {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low - 1
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
	expected := make([]int, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			arr[j], _ = strconv.Atoi(scan.Text())
		}
		expected[i] = solveCase(arr)
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
		got, err := strconv.Atoi(outScan.Text())
		if err != nil {
			fmt.Printf("bad output for test %d\n", i+1)
			os.Exit(1)
		}
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %d got %d\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
