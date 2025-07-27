package main

import "fmt"

// main is a stub solution for the quantum oracle described in problemB1.txt.
// The actual task is to implement a unitary that flips an output qubit iff
// the input register contains an equal number of zeros and ones.
// In a real quantum language one would use multi-controlled X gates
// conditioned on all bit strings with exactly N/2 ones. As Go cannot
// manipulate qubits directly, we output a constant placeholder.
func main() {
	fmt.Println(0)
}
