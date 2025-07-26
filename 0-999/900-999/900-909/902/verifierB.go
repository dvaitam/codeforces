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

func solveCase(parent, color []int) int {
	n := len(parent) - 1
	steps := 1
	for i := 2; i <= n; i++ {
		if color[i] != color[parent[i]] {
			steps++
		}
	}
	return steps
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(100) + 1
	parent := make([]int, n+1)
	for i := 2; i <= n; i++ {
		parent[i] = rng.Intn(i-1) + 1
	}
	color := make([]int, n+1)
	for i := 1; i <= n; i++ {
		color[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 2; i <= n; i++ {
		if i > 2 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", parent[i])
	}
	if n >= 2 {
		sb.WriteByte('\n')
	}
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", color[i])
	}
	sb.WriteByte('\n')
	return sb.String(), solveCase(parent, color)
}

func runCase(bin, input string, expected int) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	var got int
	if _, err := fmt.Sscan(fields[0], &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
