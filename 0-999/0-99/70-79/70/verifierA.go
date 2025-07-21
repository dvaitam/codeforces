package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const mod = 1000003

func powmod(base, exp int) int {
	res := 1
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func solve(n int) int {
	if n <= 0 {
		return 1
	}
	return powmod(3, n-1)
}

func runCase(bin string, n int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("%d\n", n))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(bytes.NewReader(out.Bytes()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	exp := solve(n)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(1001) // 0..1000
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
