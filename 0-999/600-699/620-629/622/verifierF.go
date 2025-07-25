package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 1000000007

func modPow(a, b int64) int64 {
	a %= mod
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func solveF(n, k int64) int64 {
	var sum int64
	for i := int64(1); i <= n; i++ {
		sum = (sum + modPow(i, k)) % mod
	}
	return sum
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(6)
	type Test struct{ n, k int64 }
	tests := []Test{{1, 0}, {10, 1}, {5, 2}, {100, 3}}
	for len(tests) < 100 {
		n := int64(rand.Intn(1000) + 1)
		k := int64(rand.Intn(6))
		tests = append(tests, Test{n, k})
	}
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n", t.n, t.k)
		expected := fmt.Sprintf("%d", solveF(t.n, t.k))
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed: n=%d k=%d expected %s got %s\n", i+1, t.n, t.k, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
