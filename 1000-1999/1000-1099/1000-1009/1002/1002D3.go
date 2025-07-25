package main

// MajorityOracle applies a quantum oracle for the 3-qubit majority function.
// It flips the output qubit y if and only if at least two qubits in x are in state |1>.
// The input slice x must have length 3.
func MajorityOracle(x []Qubit, y Qubit) {
   // Flip y for each pair of input qubits both in |1>
   CCNOT(x[0], x[1], y)
   CCNOT(x[0], x[2], y)
   CCNOT(x[1], x[2], y)
}
