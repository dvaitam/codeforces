package Solution

// Solve implements the quantum oracle for f(x) = x_k.
// It applies a CNOT gate from the k-th input qubit to the output qubit y,
// performing |x⟩|y⟩ → |x⟩|y ⊕ x_k⟩.
func Solve(x []Qubit, y Qubit, k int) {
    // Controlled-NOT: control is x[k], target is y
    CNOT(x[k], y)
}
