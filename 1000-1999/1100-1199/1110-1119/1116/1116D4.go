package main

import (
   "bufio"
   "fmt"
   "os"
)

// main reads N and prints an X/. pattern for the unitary described in problem D4:
// - central 2x2 sub-matrix full of X
// - anti-diagonals of top-left and bottom-right (size 2^{N-1}-1)
// - diagonals of top-right and bottom-left (size 2^{N-1}-1)
func main() {
   var N int
   if _, err := fmt.Fscan(os.Stdin, &N); err != nil {
       fmt.Fprintln(os.Stderr, "Failed to read N:", err)
       os.Exit(1)
   }
   if N < 2 || N > 5 {
       fmt.Fprintln(os.Stderr, "N out of range")
       os.Exit(1)
   }
   size := 1 << N
   half := 1 << (N - 1)
   // allocate matrix
   mat := make([][]bool, size)
   for i := range mat {
       mat[i] = make([]bool, size)
   }
   // central 2x2 block
   r0, c0 := half-1, half-1
   for dr := 0; dr < 2; dr++ {
       for dc := 0; dc < 2; dc++ {
           mat[r0+dr][c0+dc] = true
       }
   }
   // size of small submatrices
   s := half - 1
   // top-left: anti-diagonal
   for i := 0; i < s; i++ {
       mat[i][s-1-i] = true
   }
   // top-right: diagonal
   for i := 0; i < s; i++ {
       mat[i][half+1+i] = true
   }
   // bottom-left: diagonal
   for i := 0; i < s; i++ {
       mat[half+1+i][i] = true
   }
   // bottom-right: anti-diagonal
   for i := 0; i < s; i++ {
       mat[half+1+i][half+1+(s-1-i)] = true
   }
   // print pattern
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i := 0; i < size; i++ {
       for j := 0; j < size; j++ {
           if mat[i][j] {
               w.WriteByte('X')
           } else {
               w.WriteByte('.')
           }
       }
       w.WriteByte('\n')
   }
}
