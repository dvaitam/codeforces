package main

import (
   "fmt"
   "os"
   "strconv"
)

// main generates a unitary matrix on N qubits (2 <= N <= 5)
// with 2x2 non-zero blocks on the diagonal by applying an X gate
// on the least significant qubit (bit 0). The matrix is printed
// as real entries in little-endian order: for each column (input)
// indexed by bits[0]..bits[N-1] (LSB first) from 0 to 2^N-1,
// all rows from 0 to 2^N-1.
// Optionally, N can be provided as the first command-line argument.
func main() {
   // default number of qubits
   N := 2
   if len(os.Args) > 1 {
       if n, err := strconv.Atoi(os.Args[1]); err == nil && n >= 2 && n <= 5 {
           N = n
       }
   }
   size := 1 << N
   // For each column (input basis state)
   for col := 0; col < size; col++ {
       // The X gate on bit 0 flips LSB: maps |i> to |i^1>
       for row := 0; row < size; row++ {
           if row == (col ^ 1) {
               fmt.Printf("%f", 1.0)
           } else {
               fmt.Printf("%f", 0.0)
           }
           if row+1 < size || col+1 < size {
               fmt.Print(" ")
           }
       }
   }
   fmt.Println()
}
