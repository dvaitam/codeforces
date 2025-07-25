package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// This program reads the (collapsed) state of N qubits, represented as
// a binary string of length N (e.g., "0000" or "00100"), and determines
// whether the original state was |0...0> (all zeros) or the W state
// (exactly one qubit was |1>). It outputs 0 for |0...0> and 1 for W state.
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	// Read N
	if !scanner.Scan() {
		return
	}
	n, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid N")
		return
	}
	// Read state string
	if !scanner.Scan() {
		return
	}
	s := scanner.Text()
	if len(s) != n {
		fmt.Fprintln(os.Stderr, "state length mismatch")
		return
	}
	// Check if any bit is '1'
	for _, c := range s {
		if c == '1' {
			fmt.Println(1)
			return
		}
	}
	fmt.Println(0)
}
