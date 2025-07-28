package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

func generateInput() []byte {
	r := rand.New(rand.NewSource(6))
	n := 100
	var buf bytes.Buffer
	fmt.Fprintln(&buf, n)
	for i := 0; i < n; i++ {
		fmt.Fprintln(&buf, r.Intn(20)+1)
	}
	return buf.Bytes()
}

func run(bin string, args []string, input []byte) ([]byte, error) {
	cmd := exec.Command(bin, args...)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.Bytes(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	input := generateInput()

	// Run candidate binary
	candOut, err := run(bin, nil, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "execution failed: %v\n", err)
		os.Exit(1)
	}

	// Run reference solution using go run
	refPath := "1000-1999/1700-1799/1750-1759/1758/1758F.go"
	refOut, err := run("go", []string{"run", refPath}, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference execution failed: %v\n", err)
		os.Exit(1)
	}

	if !bytes.Equal(bytes.TrimSpace(candOut), bytes.TrimSpace(refOut)) {
		fmt.Fprintln(os.Stderr, "output mismatch")
		os.Exit(1)
	}
	fmt.Println("All test cases passed.")
}
