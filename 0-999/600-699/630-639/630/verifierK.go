package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod = 2520

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

var prefix [mod + 1]int

func init() {
	for i := 1; i <= mod; i++ {
		prefix[i] = prefix[i-1]
		if gcd(i, mod) == 1 {
			prefix[i]++
		}
	}
}

func solve(n int64) int64 {
	blocks := n / mod
	rem := int(n % mod)
	return int64(prefix[mod])*blocks + int64(prefix[rem])
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Int63n(1_000_000_000_000_000_000) + 1
	ans := solve(n)
	input := fmt.Sprintf("%d\n", n)
	expected := fmt.Sprintf("%d", ans)
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierK.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
