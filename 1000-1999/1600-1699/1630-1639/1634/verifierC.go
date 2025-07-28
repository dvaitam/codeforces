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

func genTests() []string {
	rng := rand.New(rand.NewSource(44))
	tests := make([]string, 100)
	for i := range tests {
		n := rng.Intn(6) + 1
		k := rng.Intn(6) + 1
		var in strings.Builder
		fmt.Fprintf(&in, "1\n%d %d\n", n, k)
		tests[i] = in.String()
	}
	return tests
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	solPath := "./solC.bin"
	if err := exec.Command("go", "build", "-o", solPath, "1634C.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(solPath)

	tests := genTests()
	for i, t := range tests {
		expect, err := run(solPath, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\nexpected:\n%s\nGot:\n%s", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
