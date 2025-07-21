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
   fmt.Fscan(reader, &n, &k)
   B := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &B[i])
   }
   // A will be built by inserting values from n down to 1
   A := make([]int, 0, n)
   for j := n; j >= 1; j-- {
       need := B[j]
       cnt := 0
       idx := 0
       // find smallest idx in [0,len(A)] such that cnt == need
       for idx = 0; idx < len(A); idx++ {
           if cnt == need {
               break
           }
           if A[idx] >= j+k {
               cnt++
           }
       }
       if idx == len(A) && cnt != need {
           // check at end
           if cnt == need {
               // ok, idx is len(A)
           } else {
               // should not happen as problem guarantees solution
           }
       }
       // insert j at position idx
       A = append(A, 0)
       copy(A[idx+1:], A[idx:])
       A[idx] = j
   }
   // output
   for i, v := range A {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
