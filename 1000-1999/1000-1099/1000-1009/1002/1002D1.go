package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

// This program implements a quantum oracle for f(x) = b·x mod 2
// on N input qubits x[1..N] and one output qubit y at position N+1.
// Input: N (number of input qubits), then a bitstring b of length N (each char '0' or '1').
// Output: a sequence of operations: for each i where b_i=='1', apply CNOT from x_i to y.
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
   // Apply CNOT from each x_i to y (qubit N+1) when b_i == '1'
   for i := 0; i < n; i++ {
       if b[i] == '1' {
           ops = append(ops, fmt.Sprintf("CNOT %d %d", i+1, n+1))
       }
   }
   // Output operations
   fmt.Println(len(ops))
   for _, op := range ops {
       fmt.Println(op)
   }
}
