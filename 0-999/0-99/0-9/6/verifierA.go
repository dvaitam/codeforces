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

func solve(sticks [4]int) string {
	triangle := false
	segment := false
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			for k := j + 1; k < 4; k++ {
				a, b, c := sticks[i], sticks[j], sticks[k]
				if a > b {
					a, b = b, a
				}
				if b > c {
					b, c = c, b
				}
				if a > b {
					a, b = b, a
				}
				if a+b > c {
					triangle = true
				} else if a+b == c {
					segment = true
				}
			}
		}
	}
	if triangle {
		return "TRIANGLE"
	}
	if segment {
		return "SEGMENT"
	}
	return "IMPOSSIBLE"
}

func generateCase(rng *rand.Rand) (string, string) {
	var sticks [4]int
	for i := 0; i < 4; i++ {
		sticks[i] = rng.Intn(100) + 1
	}
	input := fmt.Sprintf("%d %d %d %d\n", sticks[0], sticks[1], sticks[2], sticks[3])
	return input, solve(sticks)
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
		input, expected := generateCase(rng)
		if err := runCase(exe, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
