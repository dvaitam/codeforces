package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

// This program implements the Bernstein–Vazirani algorithm (E1):
// given access to an oracle f(x) = b·x mod 2 on N input qubits,
// it reconstructs the hidden bitstring b with a single oracle query.
// We output a sequence of quantum operations:
// 1. Prepare the output qubit by applying X to qubit N+1.
// 2. Apply H to all qubits 1 through N+1.
// 3. Call the oracle (black-box operation).
// 4. Apply H to input qubits 1 through N.
// The harness will then measure the input qubits to read out b.
func main() {
   scanner := bufio.NewScanner(os.Stdin)
   if !scanner.Scan() {
       return
   }
   n, err := strconv.Atoi(scanner.Text())
   if err != nil {
       fmt.Fprintln(os.Stderr, "invalid N")
       return
   }
   // Collect operations
   ops := make([]string, 0, 2*n+3)
   // Flip the output qubit (index n+1) to |1>
   ops = append(ops, fmt.Sprintf("X %d", n+1))
   // Apply Hadamard to all qubits (input and output)
   for i := 1; i <= n+1; i++ {
       ops = append(ops, fmt.Sprintf("H %d", i))
   }
   // Call the oracle
   ops = append(ops, "ORACLE")
   // Apply Hadamard to input qubits to extract b
   for i := 1; i <= n; i++ {
       ops = append(ops, fmt.Sprintf("H %d", i))
   }
   // Output number of operations and the sequence
   fmt.Println(len(ops))
   for _, op := range ops {
       fmt.Println(op)
   }
}
