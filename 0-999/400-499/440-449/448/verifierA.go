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

type testCase struct {
	input    string
	expected string
}

func expectedA(a1, a2, a3, b1, b2, b3, n int) string {
	totalCups := a1 + a2 + a3
	totalMedals := b1 + b2 + b3
	cupShelves := (totalCups + 4) / 5
	medalShelves := (totalMedals + 9) / 10
	if cupShelves+medalShelves <= n {
		return "YES"
	}
	return "NO"
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		a1 := rng.Intn(101)
		a2 := rng.Intn(101)
		a3 := rng.Intn(101)
		b1 := rng.Intn(101)
		b2 := rng.Intn(101)
		b3 := rng.Intn(101)
		n := rng.Intn(100) + 1
		input := fmt.Sprintf("%d %d %d\n%d %d %d\n%d\n", a1, a2, a3, b1, b2, b3, n)
		exp := expectedA(a1, a2, a3, b1, b2, b3, n)
		cases[i] = testCase{input: input, expected: exp}
	}
	return cases
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
	cases := generateTests()
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if out != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
