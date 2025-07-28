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

func run(bin string, input string) (string, error) {
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func expected(n int) (int, int, int) {
	a := (n + 1) / 3
	b := (n+2)/3 + 1
	c := n/3 - 1
	return a, b, c
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	cases := []int{6, 7, 8, 9, 10, 11}
	for i := 0; i < 94; i++ {
		cases = append(cases, rng.Intn(100000-6+1)+6)
	}
	for i, n := range cases {
		input := fmt.Sprintf("1\n%d\n", n)
		expA, expB, expC := expected(n)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) < 3 {
			fmt.Fprintf(os.Stderr, "case %d: expected three numbers, got %q\n", i+1, out)
			os.Exit(1)
		}
		var a, b, c int
		fmt.Sscanf(strings.Join(fields[:3], " "), "%d %d %d", &a, &b, &c)
		if a != expA || b != expB || c != expC {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d %d %d got %d %d %d\n", i+1, expA, expB, expC, a, b, c)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
