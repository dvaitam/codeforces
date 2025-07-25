package main

// This file implements a quantum oracle on N qubits that checks whether the input
// bit-vector x is a palindrome, flipping the output qubit y if and only if x is a palindrome.
// The function signature and integration with the quantum simulator/testing harness
// are assumed to be provided externally.

// PalindromeOracle applies the oracle |x⟩|y⟩ → |x⟩|y ⊕ f(x)⟩, where f(x)=1 if x is a palindrome.
// x: input register qubits, y: output qubit, anc: available ancilla qubits (at least floor(len(x)/2)).
func PalindromeOracle(x []Qubit, y Qubit, anc []Qubit) {
   n := len(x)
   // Compute parity of each symmetric pair into ancilla qubits
   for i := 0; i < n/2; i++ {
       // anc[i] = x[i] XOR x[n-1-i]
       CNOT(x[i], anc[i])
       CNOT(x[n-1-i], anc[i])
   }
   // Flip y if all ancillas are 0 (i.e., all pairs equal), i.e., x is palindrome
   if len(anc) > 1 {
       // multi-controlled X with controls on anc (checking all zeros)
       MultiControlledXInvert(anc, y)
   } else if len(anc) == 1 {
       // single control: anc[0]==0 implies equality, flip y
       X(anc[0])
       CNOT(anc[0], y)
       X(anc[0])
   } else {
       // n <= 1: always palindrome, just flip y
       X(y)
   }
   // Uncompute ancillas
   for i := n/2 - 1; i >= 0; i-- {
       CNOT(x[n-1-i], anc[i])
       CNOT(x[i], anc[i])
   }
}

func main() {
   // Entry point left empty; the testing harness should invoke PalindromeOracle.
}
