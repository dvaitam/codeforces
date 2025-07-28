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

func solve(n int) string {
	var lines []string
	if n <= 3 {
		for i := 1; i < n; i++ {
			lines = append(lines, fmt.Sprintf("1 %d", i))
		}
		lines = append(lines, fmt.Sprintf("%d %d", n, n))
	} else {
		lines = append(lines, "1 1")
		for i := 3; i < n; i++ {
			lines = append(lines, fmt.Sprintf("1 %d", i))
		}
		lines = append(lines, fmt.Sprintf("%d %d", n, n-1))
		lines = append(lines, fmt.Sprintf("%d %d", n, n))
	}
	return strings.Join(lines, "\n")
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	input := fmt.Sprintf("1\n%d\n", n)
	expect := solve(n)
	return input, expect
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		inp, exp := genCase(rng)
		got, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i+1, exp, got, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
