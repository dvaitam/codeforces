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

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // Read A and mark in S
   A := make([]int, k+1)
   S := make([]bool, 901)
   for i := 1; i <= k; i++ {
       fmt.Fscan(reader, &A[i])
       if A[i] >= 0 && A[i] < len(S) {
           S[A[i]] = true
       }
   }
   // Generate output
   p := 1
   for i := 1; i <= k; i++ {
       fmt.Fprintf(writer, "%d ", A[i])
       for j := 1; j < n; p++ {
           if p < len(S) && !S[p] {
               fmt.Fprintf(writer, "%d ", p)
               j++
           }
       }
   }
   fmt.Fprintln(writer)
