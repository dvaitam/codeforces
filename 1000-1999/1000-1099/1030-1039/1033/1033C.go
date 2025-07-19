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
   // 1-indexed arrays
   A := make([]int, n+1)
   ind := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &A[i])
       ind[A[i]] = i
   }
   win := make([]bool, n+1)
   // win[n] is false by default
   for i := n - 1; i >= 1; i-- {
       pos := ind[i]
       // try jumps to the right
       w := false
       for j := pos; j <= n; j += i {
           if A[j] > i && !win[A[j]] {
               w = true
               break
           }
       }
       if !w {
           // try jumps to the left
           for j := pos; j >= 1; j -= i {
               if A[j] > i && !win[A[j]] {
                   w = true
                   break
               }
           }
       }
       win[i] = w
   }
   // output result for each position in input order
   // win[A[i]] true -> 'A', else 'B'
   for i := 1; i <= n; i++ {
       if win[A[i]] {
           writer.WriteByte('A')
       } else {
           writer.WriteByte('B')
       }
   }
   writer.WriteByte('\n')
}
