package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   A := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &A[i])
   }

   for i := 0; i < n; i++ {
       for A[i] > 0 {
           writer.WriteByte('P')
           A[i]--
           if A[i] > 0 {
               if i == n-1 {
                   writer.WriteByte('L')
                   writer.WriteByte('R')
               } else {
                   writer.WriteByte('R')
                   writer.WriteByte('L')
               }
           }
       }
       if i != n-1 {
           writer.WriteByte('R')
       }
   }
   writer.WriteByte('\n')
}
