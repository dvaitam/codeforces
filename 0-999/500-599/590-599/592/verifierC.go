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

func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func expected(t, w, b uint64) string {
	if w > b {
		w, b = b, w
	}
	g := gcd(w, b)
	l := b / g
	var L uint64
	if l > 0 && w > t/l {
		L = t + 1
	} else {
		L = l * w
		if L > t {
			L = t + 1
		}
	}
	q := t / L
	r := t % L
	num := q*w + min(r, w-1)
	den := t
	g2 := gcd(num, den)
	return fmt.Sprintf("%d/%d", num/g2, den/g2)
}

func generateCase(rng *rand.Rand) (string, string) {
	t := uint64(rng.Int63n(1_000_000_000_000) + 1)
	w := uint64(rng.Int63n(1_000_000) + 1)
	b := uint64(rng.Int63n(1_000_000) + 1)
	input := fmt.Sprintf("%d %d %d\n", t, w, b)
	return input, expected(t, w, b)
}

func runCase(exe, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
