package main

import (
   "fmt"
   "math"
   "os"
   "strconv"
)

// main generates a unitary upper Hessenberg matrix on N qubits (2 <= N <= 4)
// by applying successive Givens rotations on adjacent rows with angle pi/4.
// The resulting matrix has zeros below the first subdiagonal and non-zeros on
// the first subdiagonal and above. It is printed as real entries in little-endian
// order: for each column (input basis state) from 0 to 2^N-1, all rows from 0 to 2^N-1.
// Optionally, N can be provided as the first command-line argument.
func main() {
   // default number of qubits
   N := 2
   if len(os.Args) > 1 {
       if n, err := strconv.Atoi(os.Args[1]); err == nil && n >= 2 && n <= 4 {
           N = n
       }
   }
   size := 1 << N
   // initialize identity matrix H
   H := make([][]float64, size)
   for i := 0; i < size; i++ {
       H[i] = make([]float64, size)
       H[i][i] = 1.0
   }
   // Givens rotation parameters: cos = sin = 1/sqrt(2)
   inv := 1.0 / math.Sqrt2
   cos := inv
   sin := inv
   // apply Givens rotations on rows (i, i+1) for i=0..size-2
   for i := 0; i < size-1; i++ {
       for col := 0; col < size; col++ {
           a := H[i][col]
           b := H[i+1][col]
           H[i][col] = cos*a - sin*b
           H[i+1][col] = sin*a + cos*b
       }
   }
   // print matrix by columns, rows in little-endian order
   for col := 0; col < size; col++ {
       for row := 0; row < size; row++ {
           fmt.Printf("%f", H[row][col])
           if row+1 < size || col+1 < size {
               fmt.Print(" ")
           }
       }
   }
   fmt.Println()
}
