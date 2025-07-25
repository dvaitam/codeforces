package main

// 1115U1: implements an anti-diagonal unitary on N qubits by applying X to each qubit.
// Input: none (qubits provided by harness)
// Output: none (state modified in place)

// AntiDiagonal applies a Pauli-X gate to each qubit in the register.
func AntiDiagonal(qubits []Qubit) {
   for i := range qubits {
       qubits[i].X()
   }
}

// Qubit represents a single quantum bit in the harness environment.
// The harness should provide a concrete implementation of X().
type Qubit interface {
   X()
}
