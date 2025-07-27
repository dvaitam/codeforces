package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(a, b int64) int64 {
	rem := a % b
	if rem == 0 {
		return 0
	}
	return b - rem
}

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	for t := 1; t <= 100; t++ {
		a := rand.Int63n(1_000_000_000) + 1
		b := rand.Int63n(1_000_000_000) + 1
		input := fmt.Sprintf("1\n%d %d\n", a, b)
		expect := fmt.Sprintf("%d", solve(a, b))
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", t, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\ninput:\n%s", t, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
