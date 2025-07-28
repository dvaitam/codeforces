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

func expected(n int) string {
	if n%2 == 0 {
		return "-1"
	}
	ops := make([]int, 0)
	for n > 1 {
		if n%4 == 1 {
			ops = append(ops, 1)
			n = (n + 1) / 2
		} else {
			ops = append(ops, 2)
			n = (n - 1) / 2
		}
	}
	// reverse ops
	for i, j := 0, len(ops)-1; i < j; i, j = i+1, j-1 {
		ops[i], ops[j] = ops[j], ops[i]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprint(len(ops)))
	if len(ops) > 0 {
		sb.WriteByte('\n')
		for i, v := range ops {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
	}
	return sb.String()
}

func runCase(exe string, n int) error {
	input := fmt.Sprintf("1\n%d\n", n)
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected(n))
	if got != exp {
		return fmt.Errorf("n=%d expected %q got %q", n, exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		var n int
		if rng.Intn(2) == 0 {
			n = rng.Intn(1_000_000_000) + 1
		} else {
			n = rng.Intn(1_000_000_000/2)*2 + 1
		}
		if err := runCase(exe, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
