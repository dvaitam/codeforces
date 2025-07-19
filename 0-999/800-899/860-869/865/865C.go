package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, R int
   if _, err := fmt.Fscan(in, &n, &R); err != nil {
       return
   }
   f := make([]int, n)
   s := make([]int, n)
   p := make([]int, n)
   d := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &f[i], &s[i], &p[i])
       d[i] = s[i] - f[i]
       R -= f[i]
   }
   // dp[i][t]: expected value from task i with accumulated extra t
   dp := make([][]float64, n+1)
   for i := range dp {
       dp[i] = make([]float64, R+1)
   }
   // check if dp[0][0] > C
   check := func(C float64) bool {
       // reset dp
       for i := 0; i <= n; i++ {
           for t := 0; t <= R; t++ {
               dp[i][t] = 0
           }
       }
       for i := n - 1; i >= 0; i-- {
           pi := float64(p[i]) / 100.0
           qi := 1 - pi
           fi := float64(f[i])
           si := float64(s[i])
           di := d[i]
           for t := 0; t <= R; t++ {
               f1 := pi * (dp[i+1][t] + fi)
               var f2 float64
               if t+di > R {
                   f2 = qi * (C + si)
               } else {
                   f2 = qi * (dp[i+1][t+di] + si)
               }
               val := f1 + f2
               if i > 0 {
                   dp[i][t] = math.Min(C, val)
               } else {
                   dp[i][t] = val
               }
           }
       }
       return dp[0][0] > C
   }
   // binary search
   l, r := 0.0, 1e12
   for it := 0; it < 80; it++ {
       mid := (l + r) / 2
       if check(mid) {
           l = mid
       } else {
           r = mid
       }
   }
   // result
   fmt.Printf("%.17f\n", l)
}
