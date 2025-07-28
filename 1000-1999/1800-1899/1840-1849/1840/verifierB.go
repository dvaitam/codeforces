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

func solve(n int64, k int64) int64 {
	if k >= 31 {
		if n+1 < 0 { // overflow not possible with int64
			return n + 1
		}
		return n + 1
	}
	pow := int64(1) << k
	if n+1 < pow {
		return n + 1
	}
	return pow
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Int63n(1_000_000_000)
	k := rng.Int63n(40)
	if rng.Intn(4) == 0 {
		k = rng.Int63n(40) + 31
	}
	input := fmt.Sprintf("1\n%d %d\n", n, k)
	expected := fmt.Sprintf("%d", solve(n, k))
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
