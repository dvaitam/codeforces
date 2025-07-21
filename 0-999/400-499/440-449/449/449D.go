package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   maxA := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       if a[i] > maxA {
           maxA = a[i]
       }
   }
   // determine bit width
   B := 0
   for (1 << B) <= maxA {
       B++
   }
   size := 1 << B
   // frequency
   f := make([]int, size)
   for _, v := range a {
       f[v]++
   }
   // superset zeta transform: f[mask] = count of values superset of mask
   for i := 0; i < B; i++ {
       bit := 1 << i
       for mask := 0; mask < size; mask++ {
           if mask&bit == 0 {
               f[mask] += f[mask|bit]
           }
       }
   }
   // precompute powers of two
   pow2 := make([]int, n+1)
   pow2[0] = 1
   for i := 1; i <= n; i++ {
       pow2[i] = pow2[i-1] * 2 % mod
   }
   // p[mask] = number of non-empty subsequences of elements superset of mask
   dp := make([]int, size)
   for mask := 0; mask < size; mask++ {
       cnt := f[mask]
       if cnt > 0 {
           dp[mask] = pow2[cnt] - 1
           if dp[mask] < 0 {
               dp[mask] += mod
           }
       }
   }
   // mobius inversion (inverse superset transform) to get exact AND == mask counts
   for i := 0; i < B; i++ {
       bit := 1 << i
       for mask := 0; mask < size; mask++ {
           if mask&bit == 0 {
               dp[mask] -= dp[mask|bit]
               if dp[mask] < 0 {
                   dp[mask] += mod
               }
           }
       }
   }
   // result for AND == 0
   fmt.Println(dp[0])
}
