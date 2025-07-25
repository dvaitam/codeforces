package main

import (
   "bufio"
   "fmt"
   "os"
)

// This program prints a unitary matrix pattern for U3:
// - top-left quarter is an anti-diagonal block of non-zero elements (X),
// - top-right and bottom-left quarters are zeros (.),
// - bottom-right quarter is filled with non-zero elements (X).
// Input: a single integer N (2 <= N <= 5).
// Output: a 2^N x 2^N pattern of 'X' and '.' characters.
func main() {
   reader := bufio.NewReader(os.Stdin)
   var N int
   if _, err := fmt.Fscan(reader, &N); err != nil {
       return
   }
   size := 1 << N
   half := size >> 1
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < size; i++ {
       for j := 0; j < size; j++ {
           if i < half {
               if j < half {
                   // top-left quarter: anti-diagonal
                   if i+j == half-1 {
                       writer.WriteByte('X')
                   } else {
                       writer.WriteByte('.')
                   }
               } else {
                   // top-right quarter: zero
                   writer.WriteByte('.')
               }
           } else {
               if j < half {
                   // bottom-left quarter: zero
                   writer.WriteByte('.')
               } else {
                   // bottom-right quarter: non-zero
                   writer.WriteByte('X')
               }
           }
       }
       writer.WriteByte('\n')
   }
}
