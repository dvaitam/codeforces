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

func genCase(rng *rand.Rand) (string, string) {
	n1 := rng.Intn(50) + 1
	n2 := rng.Intn(50) + 1
	k1 := rng.Intn(50) + 1
	k2 := rng.Intn(50) + 1
	input := fmt.Sprintf("%d %d %d %d\n", n1, n2, k1, k2)
	expected := "Second\n"
	if n1 > n2 {
		expected = "First\n"
	}
	return input, expected
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 1; i <= cases; i++ {
		input, expect := genCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\n\ngot:\n%s\n", i, strings.TrimSpace(expect), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
