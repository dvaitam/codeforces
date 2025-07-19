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

   var N, M int
   if _, err := fmt.Fscan(in, &N, &M); err != nil {
       return
   }
   const prime = 100003
   // output header values
   fmt.Fprintf(out, "%d %d\n", prime, prime)
   // build a simple path from 1 to N
   for i := 1; i < N-1; i++ {
       fmt.Fprintf(out, "%d %d %d\n", i, i+1, 1)
   }
   // last edge to make shortest path equal to prime
   fmt.Fprintf(out, "%d %d %d\n", N-1, N, prime-N+2)
   // add remaining edges with large weight
   cnt := M - N + 1
   if cnt <= 0 {
       return
   }
   for i := 1; i < N; i++ {
       for j := i + 2; j <= N; j++ {
           fmt.Fprintf(out, "%d %d %d\n", i, j, prime+1)
           cnt--
           if cnt == 0 {
               return
           }
       }
   }
}
