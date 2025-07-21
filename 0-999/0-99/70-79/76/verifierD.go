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

func solve(input string) string {
	in := strings.Fields(input)
	if len(in) < 2 {
		return "-1\n"
	}
	A, _ := strconv.ParseUint(in[0], 10, 64)
	B, _ := strconv.ParseUint(in[1], 10, 64)
	if A < B {
		return "-1\n"
	}
	diff := A - B
	if diff%2 == 1 {
		return "-1\n"
	}
	d := diff / 2
	if d&B != 0 {
		return "-1\n"
	}
	x := d
	y := d + B
	return fmt.Sprintf("%d %d\n", x, y)
}

func generateCase(rng *rand.Rand) (string, string) {
	A := rng.Uint64() % 1000
	B := rng.Uint64() % 1000
	input := fmt.Sprintf("%d\n%d\n", A, B)
	expected := solve(fmt.Sprintf("%d %d", A, B))
	return input, expected
}

func runCase(exe string, in, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expected), strings.TrimSpace(out.String()))
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
