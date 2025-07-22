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

func expected(n int) int {
	m := n - 10
	switch {
	case m < 1 || m > 11:
		return 0
	case m == 10:
		return 15
	default:
		return 4
	}
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Deterministic tests for all possible n
	for n := 1; n <= 25; n++ {
		input := fmt.Sprintf("%d\n", n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "n=%d: %v\n", n, err)
			os.Exit(1)
		}
		exp := fmt.Sprintf("%d", expected(n))
		if got != exp {
			fmt.Fprintf(os.Stderr, "n=%d: expected %s got %s\n", n, exp, got)
			os.Exit(1)
		}
	}

	// Additional random tests
	for i := 0; i < 75; i++ {
		n := rng.Intn(25) + 1
		input := fmt.Sprintf("%d\n", n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "random case %d (n=%d): %v\n", i+1, n, err)
			os.Exit(1)
		}
		exp := fmt.Sprintf("%d", expected(n))
		if got != exp {
			fmt.Fprintf(os.Stderr, "random case %d (n=%d): expected %s got %s\n", i+1, n, exp, got)
			os.Exit(1)
		}
	}

	// Repeat to make at least 100 cases
	for n := 1; n <= 25; n++ {
		input := fmt.Sprintf("%d\n", n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "repeat n=%d: %v\n", n, err)
			os.Exit(1)
		}
		exp := fmt.Sprintf("%d", expected(n))
		if got != exp {
			fmt.Fprintf(os.Stderr, "repeat n=%d: expected %s got %s\n", n, exp, got)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}
