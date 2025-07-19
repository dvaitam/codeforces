package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var N int
   if _, err := fmt.Fscan(in, &N); err != nil {
       return
   }
   A := make([]int64, N)
   B := make([]int64, N)
   var sumA int64
   for i := 0; i < N; i++ {
       fmt.Fscan(in, &A[i], &B[i])
       sumA += A[i]
   }
   // Greedy assignment to keep difference within 500
   totalA := sumA
   totalB := int64(0)
   res := make([]byte, N)
   for i := range res {
       res[i] = 'A'
   }
   for i := 0; i < N; i++ {
       if totalA > totalB {
           if totalA-totalB > 500 {
               totalA -= A[i]
               totalB += B[i]
               res[i] = 'G'
           }
       } else {
           if totalB-totalA > 500 {
               fmt.Fprintln(out, -1)
               return
           }
       }
   }
   out.Write(res)
   out.WriteByte('\n')
}
