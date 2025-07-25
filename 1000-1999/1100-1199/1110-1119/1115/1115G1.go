package main

// 1115G1: Quantum oracle for the AND function f(x) = x0 ∧ x1 ∧ … ∧ xN-1.
// The operation acts on an input register x of N qubits and an output qubit y,
// performing |x⟩|y⟩ → |x⟩|y ⊕ f(x)⟩.
// AndOracle flips y if and only if all qubits in x are 1.
// No classical I/O; the testing harness should invoke AndOracle directly.
func AndOracle(x []Qubit, y Qubit) {
    // Apply a multi-controlled X gate: flip y when all controls x[i] are |1⟩.
    MultiControlledX(x, y)
}

func main() {
    // Entry point left empty; the harness will call AndOracle.
}
