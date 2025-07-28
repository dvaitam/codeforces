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
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func solveCase(a, b []int) string {
	diff := -1
	for i := range a {
		if b[i] > a[i] {
			return "NO"
		}
		if b[i] > 0 {
			d := a[i] - b[i]
			if diff == -1 {
				diff = d
			} else if diff != d {
				return "NO"
			}
		}
	}
	if diff != -1 {
		for i := range a {
			if b[i] == 0 && a[i] > diff {
				return "NO"
			}
		}
	}
	return "YES"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(43))
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(10) + 1
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(20)
		}
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				b[i] = 0
			} else {
				val := rng.Intn(a[i] + 1)
				b[i] = val
			}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", a[i])
		}
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", b[i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveCase(a, b)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", tc+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", tc+1, expected, strings.TrimSpace(out), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
