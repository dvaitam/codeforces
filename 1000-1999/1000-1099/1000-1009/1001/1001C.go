package Solution

// Solve creates a GHZ state on the given qubits:
// |GHZ> = (|00...0> + |11...1>)/sqrt(2).
// It applies a Hadamard gate on the first qubit,
// then CNOT gates from the first qubit to each other qubit.
func Solve(qs []Qubit) {
    if len(qs) == 0 {
        return
    }
    H(qs[0])
    for i := 1; i < len(qs); i++ {
        CNOT(qs[0], qs[i])
    }
}
