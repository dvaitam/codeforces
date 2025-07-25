package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func isPrime(n uint64) bool {
	if n < 2 {
		return false
	}
	small := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}
	for _, p := range small {
		if n%p == 0 {
			return n == p
		}
	}
	d := n - 1
	s := 0
	for d&1 == 0 {
		d >>= 1
		s++
	}
	bases := []uint64{2, 7, 61}
	for _, a := range bases {
		if a >= n {
			continue
		}
		if !millerTest(a, d, s, n) {
			return false
		}
	}
	return true
}

func millerTest(a, d uint64, s int, n uint64) bool {
	x := modPow(a, d, n)
	if x == 1 || x == n-1 {
		return true
	}
	for i := 1; i < s; i++ {
		x = mulMod(x, x, n)
		if x == n-1 {
			return true
		}
	}
	return false
}

func modPow(a, e, m uint64) uint64 {
	res := uint64(1)
	a %= m
	for e > 0 {
		if e&1 == 1 {
			res = mulMod(res, a, m)
		}
		a = mulMod(a, a, m)
		e >>= 1
	}
	return res
}

func mulMod(a, b, mod uint64) uint64 {
	return (a * b) % mod
}

func solveD(n uint64) int {
	if isPrime(n) {
		return 1
	}
	if n%2 == 0 {
		return 2
	}
	if isPrime(n - 2) {
		return 2
	}
	return 3
}

func genTestD() uint64 {
	return uint64(rand.Int63n(2_000_000_000-1)) + 2
}

func runBinary(path string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	path := os.Args[1]
	for i := 0; i < 100; i++ {
		n := genTestD()
		input := fmt.Sprintf("%d\n", n)
		expected := solveD(n)
		gotStr, err := runBinary(path, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		var got int
		_, err = fmt.Sscanf(gotStr, "%d", &got)
		if err != nil {
			fmt.Printf("test %d: parse output error: %v\ninput:%soutput:%s\n", i+1, err, input, gotStr)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\ninput:%sexpected:%d\ngot:%d\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
