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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   a := make([]int, t)
   b := make([]int, t)
   maxA := 0
   for i := 0; i < t; i++ {
       fmt.Fscan(reader, &a[i], &b[i])
       if a[i] > maxA {
           maxA = a[i]
       }
   }
   // Sieve for smallest prime factor up to maxA
   spf := make([]int, maxA+1)
   for i := 2; i <= maxA; i++ {
       if spf[i] == 0 {
           for j := i; j <= maxA; j += i {
               if spf[j] == 0 {
                   spf[j] = i
               }
           }
       }
   }
   // Prefix sums of bigOmega (total prime factors with multiplicity)
   F := make([]int64, maxA+1)
   // F[0] = F[1] = 0 by default
   for i := 2; i <= maxA; i++ {
       p := spf[i]
       prev := i / p
       // f(prev) = F[prev] - F[prev-1]
       fprev := F[prev] - F[prev-1]
       fi := fprev + 1
       F[i] = F[i-1] + fi
   }
   // Answer queries
   for i := 0; i < t; i++ {
       res := F[a[i]] - F[b[i]]
       writer.WriteString(fmt.Sprintf("%d\n", res))
   }
}
