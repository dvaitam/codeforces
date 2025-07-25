package main

import (
   "bufio"
   "fmt"
   "os"
)

// Generates and prints the pattern for the unitary described in problem D2:
// A 2^N by 2^N matrix where the top-right and bottom-left quarters are zero,
// the top-left quarter is the same pattern for N-1, and the bottom-right quarter is all non-zero (X).
func main() {
   // Read N from stdin
   var N int
   if _, err := fmt.Fscan(os.Stdin, &N); err != nil {
       fmt.Fprintln(os.Stderr, "Failed to read N:", err)
       os.Exit(1)
   }
   size := 1 << N
   // Initialize matrix of false (.)
   mat := make([][]bool, size)
   for i := range mat {
       mat[i] = make([]bool, size)
   }
   // Fill recursively
   var fill func(n, r0, c0 int)
   fill = func(n, r0, c0 int) {
       if n == 1 {
           // 2x2 base: X . / . X
           mat[r0][c0] = true
           mat[r0+1][c0+1] = true
           return
       }
       half := 1 << (n - 1)
       // Top-left: recurse
       fill(n-1, r0, c0)
       // Bottom-right: fill all X
       for i := r0 + half; i < r0+2*half; i++ {
           for j := c0 + half; j < c0+2*half; j++ {
               mat[i][j] = true
           }
       }
       // Other quarters (top-right, bottom-left) remain false
   }
   fill(N, 0, 0)
   // Print pattern
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < size; i++ {
       for j := 0; j < size; j++ {
           if mat[i][j] {
               writer.WriteByte('X')
           } else {
               writer.WriteByte('.')
           }
       }
       writer.WriteByte('\n')
   }
}
