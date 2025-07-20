package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

const mod int64 = 1e9 + 7

func modPow(a, b, m int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % m
		}
		a = a * a % m
		b >>= 1
	}
	return res
}

func prepareFact(n int) ([]int64, []int64) {
	fact := make([]int64, n+1)
	inv := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	inv[n] = modPow(fact[n], mod-2, mod)
	for i := n; i > 0; i-- {
		inv[i-1] = inv[i] * int64(i) % mod
	}
	return fact, inv
}

func C(n, r int64, fact, inv []int64) int64 {
	if r < 0 || r > n {
		return 0
	}
	return fact[n] * inv[r] % mod * inv[n-r] % mod
}

func solve(n int, fact, inv []int64) int64 {
	res := int64(0)
	half := (n - 1) / 2
	for s := 0; s <= half; s++ {
		for x := s + 1; x <= 2*s+1 && x <= n; x++ {
			ways := C(int64(x-1), int64(s), fact, inv) * C(int64(n-x), int64(n-2*s-1), fact, inv) % mod
			res = (res + int64(x)*ways) % mod
		}
	}
	for s := (n + 1) / 2; s <= n; s++ {
		ways := C(int64(n), int64(s), fact, inv)
		mex := int64(2*s + 1)
		res = (res + mex*ways) % mod
	}
	return res % mod
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("problemG.txt")
	if err != nil {
		fmt.Println("could not read problemG.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	maxN := 5000
	fact, inv := prepareFact(maxN)
	expected := make([]int64, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		expected[i] = solve(n, fact, inv)
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
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got, _ := strconv.ParseInt(outScan.Text(), 10, 64)
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
