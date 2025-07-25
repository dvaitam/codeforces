package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

// This program implements a quantum oracle for f(x) = (b·x + (1-b)·(1-x)) mod 2
// on N input qubits x[1..N] and one output qubit y at position N+1.
// Input: N (number of qubits in x), then a bitstring b of length N (each char '0' or '1').
// Output: a sequence of operations: first apply X to each x_i where b_i=='0',
// then CNOT from each x_i to y, then X again to those x_i to restore.
// Qubits are 1-indexed; y is qubit N+1.
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
   if !scanner.Scan() {
       return
   }
   b := scanner.Text()
   if len(b) != n {
       fmt.Fprintln(os.Stderr, "bitstring length mismatch")
       return
   }
   var ops []string
   // Flip inputs with b_i == '0'
   for i := 0; i < n; i++ {
       if b[i] == '0' {
           ops = append(ops, fmt.Sprintf("X %d", i+1))
       }
   }
   // Apply CNOT from each x_i to y
   for i := 1; i <= n; i++ {
       ops = append(ops, fmt.Sprintf("CNOT %d %d", i, n+1))
   }
   // Restore inputs
   for i := 0; i < n; i++ {
       if b[i] == '0' {
           ops = append(ops, fmt.Sprintf("X %d", i+1))
       }
   }
   // Output operations
   fmt.Println(len(ops))
   for _, op := range ops {
       fmt.Println(op)
   }
}
