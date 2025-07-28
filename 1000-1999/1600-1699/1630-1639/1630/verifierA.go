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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func generateCase(rng *rand.Rand) (string, int, int) {
	// Occasionally generate the impossible case n=4,k=3
	if rng.Intn(10) == 0 {
		return "1\n4 3\n", 4, 3
	}
	pow := rng.Intn(7) + 2 // 2..8 -> n up to 256
	n := 1 << pow
	k := rng.Intn(n)
	return fmt.Sprintf("1\n%d %d\n", n, k), n, k
}

func check(n, k int, output string) error {
	output = strings.TrimSpace(output)
	if output == "-1" {
		if n == 4 && k == n-1 {
			return nil
		}
		return fmt.Errorf("unexpected -1")
	}
	lines := strings.Split(output, "\n")
	if len(lines) != n/2 {
		return fmt.Errorf("expected %d lines got %d", n/2, len(lines))
	}
	used := make([]bool, n)
	sum := 0
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line %q", line)
		}
		var a, b int
		if _, err := fmt.Sscan(parts[0], &a); err != nil {
			return fmt.Errorf("invalid number in line %q", line)
		}
		if _, err := fmt.Sscan(parts[1], &b); err != nil {
			return fmt.Errorf("invalid number in line %q", line)
		}
		if a < 0 || a >= n || b < 0 || b >= n {
			return fmt.Errorf("numbers out of range in line %q", line)
		}
		if used[a] || used[b] {
			return fmt.Errorf("element repeated in line %q", line)
		}
		used[a] = true
		used[b] = true
		sum += a & b
	}
	for i, u := range used {
		if !u {
			return fmt.Errorf("element %d missing", i)
		}
	}
	if sum != k {
		return fmt.Errorf("AND sum mismatch: expected %d got %d", k, sum)
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
		input, n, k := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := check(n, k, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
