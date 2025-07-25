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

func expectedOutput(n int) string {
	if n <= 2 {
		var b strings.Builder
		fmt.Fprintln(&b, 1)
		for i := 0; i < n; i++ {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteByte('1')
		}
		if n > 0 {
			b.WriteByte('\n')
		}
		return strings.TrimSpace(b.String())
	}
	// sieve up to n+1
	limit := n + 1
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	var b strings.Builder
	fmt.Fprintln(&b, 2)
	for i := 2; i <= limit; i++ {
		if i > 2 {
			b.WriteByte(' ')
		}
		if isPrime[i] {
			b.WriteByte('1')
		} else {
			b.WriteByte('2')
		}
	}
	b.WriteByte('\n')
	return strings.TrimSpace(b.String())
}

func run(binary, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 1; i <= 100; i++ {
		n := i
		inp := fmt.Sprintf("%d\n", n)
		exp := expectedOutput(n)
		got, err := run(bin, inp)
		if err != nil {
			fmt.Printf("Test %d: execution error: %v\nOutput:\n%s\n", i, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("Test %d failed. Input:%d Expected:%s Got:%s\n", i, n, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
