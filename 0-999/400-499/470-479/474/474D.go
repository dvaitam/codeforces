package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var t, k int
   if _, err := fmt.Fscan(in, &t, &k); err != nil {
       return
   }
   pairs := make([][2]int, t)
   maxb := 0
   for i := 0; i < t; i++ {
       a, b := 0, 0
       fmt.Fscan(in, &a, &b)
       pairs[i][0] = a
       pairs[i][1] = b
       if b > maxb {
           maxb = b
       }
   }
   const mod = 1000000007
   dp := make([]int, maxb+1)
   dp[0] = 1
   for i := 1; i <= maxb; i++ {
       dp[i] = dp[i-1]
       if i >= k {
           dp[i] += dp[i-k]
           if dp[i] >= mod {
               dp[i] -= mod
           }
       }
   }
   pref := make([]int, maxb+1)
   for i := 1; i <= maxb; i++ {
       pref[i] = pref[i-1] + dp[i]
       if pref[i] >= mod {
           pref[i] -= mod
       }
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i := 0; i < t; i++ {
       a := pairs[i][0]
       b := pairs[i][1]
       res := pref[b]
       if a > 1 {
           res -= pref[a-1]
       }
       if res < 0 {
           res += mod
       }
       fmt.Fprintln(out, res)
   }
}
