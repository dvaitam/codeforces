package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func hasInstruction(p string) bool {
	return strings.ContainsAny(p, "HQ9")
}

func randomString(r *rand.Rand) string {
	n := r.Intn(100) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(r.Intn(94) + 33) // printable ASCII
	}
	return string(b)
}

func run(binary, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(44))
	for i := 1; i <= 100; i++ {
		s := randomString(r)
		input := fmt.Sprintf("%s\n", s)
		expected := "NO"
		if hasInstruction(s) {
			expected = "YES"
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Printf("wrong answer on test %d: expected %s got %s\n", i, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
