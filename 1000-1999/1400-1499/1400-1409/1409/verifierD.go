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

func digitSum(n uint64) uint64 {
	var s uint64
	for n > 0 {
		s += n % 10
		n /= 10
	}
	return s
}

func solve(n, s uint64) uint64 {
	if digitSum(n) <= s {
		return 0
	}
	var add uint64
	var pow10 uint64 = 1
	for i := 0; i < 19; i++ {
		digit := (n / pow10) % 10
		delta := (10 - digit) * pow10
		if delta > 0 {
			add += delta
			n += delta
		}
		if digitSum(n) <= s {
			break
		}
		pow10 *= 10
	}
	return add
}

func generateCase(rng *rand.Rand) (string, uint64) {
	n := uint64(rng.Int63n(1_000_000_000_000)) + 1
	s := uint64(rng.Intn(162) + 1)
	input := fmt.Sprintf("1\n%d %d\n", n, s)
	expected := solve(n, s)
	return input, expected
}

func runCase(exe, input string, expected uint64) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	var got uint64
	if _, err := fmt.Sscan(outStr, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
