package main

import (
	"bufio"
	"fmt"
	"os"
)

// generateOps builds the sequence of gates needed for the target superposition.
func generateOps(n int, bits string) ([]string, error) {
	if len(bits) != n {
		return nil, fmt.Errorf("bitstring length mismatch: got %d want %d", len(bits), n)
	}
	ops := []string{"H 1"}
	for i := 1; i < n; i++ {
		if bits[i] == '1' {
			ops = append(ops, fmt.Sprintf("CNOT 1 %d", i+1))
		}
	}
	return ops, nil
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		fmt.Fprintln(os.Stderr, "invalid N")
		return
	}
	var bits string
	if _, err := fmt.Fscan(in, &bits); err != nil {
		fmt.Fprintln(os.Stderr, "missing bitstring")
		return
	}
	ops, err := generateOps(n, bits)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(len(ops))
	for _, op := range ops {
		fmt.Println(op)
	}
}
