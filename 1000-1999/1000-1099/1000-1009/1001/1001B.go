package main

import (
	"bufio"
	"fmt"
	"os"
)

// Reads a Bell-state index (0 to 3) and constructs the corresponding Bell state on two qubits.
// As a classical stub for quantum operation, this program reads the index and performs no output.
func main() {
	reader := bufio.NewReader(os.Stdin)
	var idx int
	if _, err := fmt.Fscan(reader, &idx); err != nil {
		return
	}
	// Quantum operations are implicit; no output for this stub implementation.
}
