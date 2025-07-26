package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expected(w, h, k int) int {
	res := 0
	for i := 0; i < k; i++ {
		res += 2*h + 2*w - 4
		w -= 4
		h -= 4
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	w := rng.Intn(98) + 3
	h := rng.Intn(98) + 3
	maxK := (min(w, h) + 1) / 4
	if maxK < 1 {
		maxK = 1
	}
	k := rng.Intn(maxK) + 1
	input := fmt.Sprintf("%d %d %d\n", w, h, k)
	out := fmt.Sprintf("%d\n", expected(w, h, k))
	return input, out
}

func runCase(bin, in, out string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, errBuf.String())
	}
	got := strings.TrimSpace(buf.String())
	exp := strings.TrimSpace(out)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		in, out := generateCase(rng)
		if err := runCase(bin, in, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
