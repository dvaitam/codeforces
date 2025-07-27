package main

// 1356A3: Distinguish between the Z gate and the S gate using two applications of
// an unknown single-qubit operation. In a quantum environment one would:
//  1. Prepare a qubit in the |+> state using a Hadamard.
//  2. Apply the given operation twice.
//  3. Apply a Hadamard and measure in the computational basis.
//
// A result of 0 (state |+>) indicates the operation was Z; a result of 1 (state
// |->) indicates the operation was S.
//
// Quantum operations are not available in this Go repository. Therefore this file
// only documents the intended algorithm and leaves main empty.
func main() {
	// Implementation would require quantum gate primitives.
}
