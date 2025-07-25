package main

// GenerateSuperposition applies a Hadamard gate to each qubit, resulting
// in an equal superposition of all 2^N basis states for the N-qubit input.
// This operation takes an array of N qubits and has no return value;
// its effect is to modify the state of the qubits in place.
func GenerateSuperposition(qs []Qubit) {
   for i := range qs {
       H(qs[i])
   }
}
