package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expectedA(n int64) int64 {
	total := n * (n + 1) / 2
	var sum int64
	for p := int64(1); p <= n; p <<= 1 {
		sum += p
	}
	return total - 2*sum
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	var out strings.Builder
	sb.WriteString(strconv.Itoa(t) + "\n")
	for i := 0; i < t; i++ {
		n := rng.Int63n(1_000_000_000) + 1
		sb.WriteString(strconv.FormatInt(n, 10))
		sb.WriteByte('\n')
		out.WriteString(strconv.FormatInt(expectedA(n), 10))
		out.WriteByte('\n')
	}
	return sb.String(), out.String()
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected\n%s\ngot\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
