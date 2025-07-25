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

type Matrix [2][2]uint64

func mul(a, b Matrix, mod uint64) Matrix {
	return Matrix{{(a[0][0]*b[0][0] + a[0][1]*b[1][0]) % mod, (a[0][0]*b[0][1] + a[0][1]*b[1][1]) % mod},
		{(a[1][0]*b[0][0] + a[1][1]*b[1][0]) % mod, (a[1][0]*b[0][1] + a[1][1]*b[1][1]) % mod}}
}

func matrixPow(mat Matrix, exp, mod uint64) Matrix {
	res := Matrix{{1, 0}, {0, 1}}
	for exp > 0 {
		if exp&1 == 1 {
			res = mul(res, mat, mod)
		}
		mat = mul(mat, mat, mod)
		exp >>= 1
	}
	return res
}

func fib(n, mod uint64) uint64 {
	if n == 0 {
		return 0
	}
	M := Matrix{{1, 1}, {1, 0}}
	P := matrixPow(M, n-1, mod)
	return P[0][0]
}

func powMod(base, exp, mod uint64) uint64 {
	res := uint64(1) % mod
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			res = (res * base) % mod
		}
		base = (base * base) % mod
		exp >>= 1
	}
	return res
}

func expected(n, k, l, m uint64) uint64 {
	if l < 64 && k>>l != 0 {
		return 0
	}
	c1 := uint64(0)
	kk := k
	for kk > 0 {
		c1 += kk & 1
		kk >>= 1
	}
	c0 := l - c1
	A := fib(n+2, m)
	T := powMod(2, n, m)
	var B uint64
	if T >= A {
		B = T - A
	} else {
		B = T + m - A
	}
	r0 := powMod(A, c0, m)
	r1 := powMod(B, c1, m)
	return (r0 * r1) % m
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 4 {
			fmt.Printf("case %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.ParseUint(parts[0], 10, 64)
		k, _ := strconv.ParseUint(parts[1], 10, 64)
		l, _ := strconv.ParseUint(parts[2], 10, 64)
		m, _ := strconv.ParseUint(parts[3], 10, 64)
		exp := expected(n, k, l, m)
		input := fmt.Sprintf("%d %d %d %d\n", n, k, l, m)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err := strconv.ParseUint(gotStr, 10, 64)
		if err != nil || got != exp {
			fmt.Printf("case %d failed: expected %d got %s\n", idx, exp, gotStr)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
