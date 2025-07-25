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

func expected(v1, v2, t, d int) int {
	total := 0
	for i := 0; i < t; i++ {
		maxFromStart := v1 + i*d
		maxFromEnd := v2 + (t-1-i)*d
		if maxFromStart < maxFromEnd {
			total += maxFromStart
		} else {
			total += maxFromEnd
		}
	}
	return total
}

func check(v1, v2, t, d int, out string) error {
	val, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return fmt.Errorf("cannot parse output")
	}
	exp := expected(v1, v2, t, d)
	if val != exp {
		return fmt.Errorf("expected %d got %d", exp, val)
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
		v1 := rng.Intn(100) + 1
		v2 := rng.Intn(100) + 1
		t := rng.Intn(99) + 2
		d := rng.Intn(11)
		input := fmt.Sprintf("%d %d\n%d %d\n", v1, v2, t, d)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if err := check(v1, v2, t, d, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%soutput:%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
