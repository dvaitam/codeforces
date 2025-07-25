package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

// This program prepares a quantum state which is an equal superposition of |0...0> and a given basis state.
// Input: N (number of qubits), then a bitstring of length N (first bit is '1') describing the basis state |ψ>.
// Output: sequence of quantum gates to apply:
// First apply H on qubit 1, then for each i>1 where bit[i]=='1', apply CNOT from qubit 1 to qubit i.
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
   s := scanner.Text()
   if len(s) != n {
       fmt.Fprintln(os.Stderr, "state length mismatch")
       return
   }
   // Collect operations
   ops := make([]string, 0, n)
   // Apply Hadamard on first qubit
   ops = append(ops, fmt.Sprintf("H 1"))
   // For each qubit i>1 where bit is '1', apply CNOT(1, i)
   for i := 1; i < n; i++ {
       if s[i] == '1' {
           ops = append(ops, fmt.Sprintf("CNOT 1 %d", i+1))
       }
   }
   // Output number of operations and the operations
   fmt.Println(len(ops))
   for _, op := range ops {
       fmt.Println(op)
   }
}
