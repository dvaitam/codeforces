package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod = 1000000007

func phi(n int64) int64 {
	res := n
	m := n
	for i := int64(2); i*i <= m; i++ {
		if m%i == 0 {
			for m%i == 0 {
				m /= i
			}
			res = res / i * (i - 1)
		}
	}
	if m > 1 {
		res = res / m * (m - 1)
	}
	return res % mod
}

func run(bin, in string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 1; i <= 100; i++ {
		n := int64(i*100 + 1)
		k := int64(i%10 + 1)
		inp := fmt.Sprintf("%d %d\n", n, k)
		exp := fmt.Sprintf("%d", phi(n))
		got, err := run(bin, inp)
		if err != nil {
			fmt.Printf("Test %d: execution error: %v\nOutput:\n%s\n", i, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("Test %d failed. Input:%sExpected:%s Got:%s\n", i, inp, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
