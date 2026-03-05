package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func validate(n int, output string) string {
	S := int64(n) * int64(n+1) / 2
	minDiff := S % 2

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		return "expected 2 lines of output"
	}
	var diff int64
	if _, err := fmt.Sscan(strings.TrimSpace(lines[0]), &diff); err != nil {
		return fmt.Sprintf("cannot parse diff: %v", err)
	}
	if diff != minDiff {
		return fmt.Sprintf("wrong diff: got %d, want %d", diff, minDiff)
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return "second line too short"
	}
	var size int
	if _, err := fmt.Sscan(fields[0], &size); err != nil {
		return fmt.Sprintf("cannot parse size: %v", err)
	}
	if size < 1 || size >= n {
		return fmt.Sprintf("size %d out of range [1, %d]", size, n-1)
	}
	if len(fields) != size+1 {
		return fmt.Sprintf("expected %d integers after size, got %d", size, len(fields)-1)
	}

	seen := make([]bool, n+1)
	var sum1 int64
	for _, f := range fields[1:] {
		var v int
		if _, err := fmt.Sscan(f, &v); err != nil || v < 1 || v > n {
			return fmt.Sprintf("value %q out of range [1, %d]", f, n)
		}
		if seen[v] {
			return fmt.Sprintf("duplicate value %d", v)
		}
		seen[v] = true
		sum1 += int64(v)
	}
	sum2 := S - sum1
	got := sum1 - sum2
	if got < 0 {
		got = -got
	}
	if got != diff {
		return fmt.Sprintf("actual |sum1-sum2|=%d but reported diff=%d", got, diff)
	}
	return ""
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(8993))
	inputs := []string{"4\n", "3\n", "5\n", "2\n"}
	for len(inputs) < 100 {
		n := rng.Intn(60) + 2
		inputs = append(inputs, fmt.Sprintf("%d\n", n))
	}
	for i, input := range inputs {
		var n int
		fmt.Sscan(input, &n)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if msg := validate(n, out); msg != "" {
			fmt.Fprintf(os.Stderr, "test %d failed: %s\ninput:\n%soutput:\n%s\n", i+1, msg, input, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(inputs))
}
