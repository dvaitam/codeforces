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

func expected(n, x, a, b int) int {
	if a > b {
		a, b = b, a
	}
	dist := b - a
	ans := dist + x
	if ans > n-1 {
		ans = n - 1
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, []string) {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	exp := make([]string, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(100) + 1
		x := rng.Intn(100)
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, x, a, b))
		exp[i] = fmt.Sprintf("%d", expected(n, x, a, b))
	}
	return sb.String(), exp
}

func runCase(bin, input string, exp []string) error {
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
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(lines))
	}
	for i, line := range lines {
		if strings.TrimSpace(line) != exp[i] {
			return fmt.Errorf("line %d expected %s got %s", i+1, exp[i], strings.TrimSpace(line))
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
