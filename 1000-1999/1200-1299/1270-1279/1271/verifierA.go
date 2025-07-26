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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expected(a, b, c, d, e, f int) string {
	res := 0
	if e > f {
		t1 := min(a, d)
		res += t1 * e
		d -= t1
		t2 := min(min(b, c), d)
		res += t2 * f
	} else {
		t2 := min(min(b, c), d)
		res += t2 * f
		d -= t2
		t1 := min(a, d)
		res += t1 * e
	}
	return fmt.Sprintf("%d", res)
}

func generateCase(rng *rand.Rand) (string, string) {
	a := rng.Intn(100000) + 1
	b := rng.Intn(100000) + 1
	c := rng.Intn(100000) + 1
	d := rng.Intn(100000) + 1
	e := rng.Intn(1000) + 1
	f := rng.Intn(1000) + 1
	input := fmt.Sprintf("%d\n%d\n%d\n%d\n%d\n%d\n", a, b, c, d, e, f)
	return input, expected(a, b, c, d, e, f)
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
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		if err := runCase(exe, input, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
