package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var M int
   if _, err := fmt.Fscan(in, &M); err != nil {
       return
   }
   X := make([]int, M)
   Y := make([]int, M)
   N := 0
   for i := 0; i < M; i++ {
       fmt.Fscan(in, &X[i])
       N += X[i]
   }
   for i := 0; i < M; i++ {
       fmt.Fscan(in, &Y[i])
   }
   // Precompute binomial coefficients up to N
   c := make([][]int, N+1)
   for i := 0; i <= N; i++ {
       c[i] = make([]int, i+1)
       c[i][0] = 1
       for j := 1; j < i; j++ {
           c[i][j] = c[i-1][j-1] + c[i-1][j]
           if c[i][j] >= mod {
               c[i][j] -= mod
           }
       }
       if i > 0 {
           c[i][i] = 1
       }
   }
   // dp arrays
   dpPrev := make([]int, N+1)
   dpCurr := make([]int, N+1)
   dpPrev[0] = 1
   for idx := 0; idx < M; idx++ {
       xi := X[idx]
       yi := Y[idx]
       // reset dpCurr
       for i := 0; i <= N; i++ {
           dpCurr[i] = 0
       }
       for open := 0; open <= N; open++ {
           v := dpPrev[open]
           if v == 0 {
               continue
           }
           t := open + xi
           maxK := yi
           if maxK > t {
               maxK = t
           }
           // choose k groups to end here
           for k := 0; k <= maxK; k++ {
               ways := c[t][k]
               val := v * ways % mod
               dpCurr[t-k] = (dpCurr[t-k] + val) % mod
           }
       }
       // swap
       dpPrev, dpCurr = dpCurr, dpPrev
   }
   // result is dpPrev[0]
   fmt.Println(dpPrev[0])
}
