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

func solve(n, t int) string {
	if n == 1 && t == 10 {
		return "-1"
	}
	if t == 10 {
		return "1" + strings.Repeat("0", n-1)
	}
	return strings.Repeat(fmt.Sprint(t), n)
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
		expect := solve(tc.n, tc.t)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, input)
			os.Exit(1)
		}
		if out != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", idx+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
