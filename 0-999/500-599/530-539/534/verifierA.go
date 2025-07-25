package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("time limit")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func checkCase(n int, out string) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("cannot parse k")
	}
	fields = fields[1:]
	if len(fields) != k {
		return fmt.Errorf("expected %d numbers, got %d", k, len(fields))
	}
	seq := make([]int, k)
	used := make(map[int]bool)
	for i := 0; i < k; i++ {
		v, err := strconv.Atoi(fields[i])
		if err != nil {
			return fmt.Errorf("bad number")
		}
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range", v)
		}
		if used[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		used[v] = true
		seq[i] = v
	}
	for i := 0; i+1 < k; i++ {
		if abs(seq[i]-seq[i+1]) == 1 {
			return fmt.Errorf("adjacent numbers %d and %d differ by 1", seq[i], seq[i+1])
		}
	}
	expected := n
	if n == 1 || n == 2 {
		expected = 1
	} else if n == 3 {
		expected = 2
	}
	if k != expected {
		return fmt.Errorf("expected k=%d got %d", expected, k)
	}
	if n >= 4 && len(used) != n {
		return fmt.Errorf("not all students used")
	}
	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(5000) + 1
		input := fmt.Sprintf("%d\n", n)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if err := checkCase(n, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%soutput:%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
