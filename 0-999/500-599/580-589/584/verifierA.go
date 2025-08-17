package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("time limit")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func check(n, t int, out string) error {
	if n == 1 && t == 10 {
		if out == "-1" {
			return nil
		}
		return fmt.Errorf("expected -1 got %s", out)
	}
	if out == "-1" {
		return fmt.Errorf("unexpected -1")
	}
	if len(out) != n {
		return fmt.Errorf("wrong length: got %d want %d", len(out), n)
	}
	if out[0] == '0' {
		return fmt.Errorf("leading zero")
	}
	mod := 0
	for _, ch := range out {
		if ch < '0' || ch > '9' {
			return fmt.Errorf("invalid character %q", ch)
		}
		mod = (mod*10 + int(ch-'0')) % t
	}
	if mod != 0 {
		return fmt.Errorf("not divisible by %d", t)
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

	// predetermined edge cases
	cases := []struct{ n, t int }{
		{1, 10}, {1, 2}, {2, 10}, {3, 10}, {10, 3},
	}
	for i := 0; i < 95; i++ {
		n := rng.Intn(100) + 1
		t := rng.Intn(9) + 2 // 2..10
		cases = append(cases, struct{ n, t int }{n, t})
	}
	for idx, tc := range cases {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.t)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, input)
			os.Exit(1)
		}
		if err := check(tc.n, tc.t, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
