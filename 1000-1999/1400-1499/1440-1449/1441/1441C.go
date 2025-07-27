package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // dp[j]: max sum with j picks
   const inf = int64(-9e18)
   dp := make([]int64, k+1)
   for i := 1; i <= k; i++ {
       dp[i] = inf
   }
   for i := 0; i < n; i++ {
       var ti int
       fmt.Fscan(reader, &ti)
       u := ti
       if u > k {
           u = k
       }
       // read prefix of length u
       prefix := make([]int64, u+1)
       for j := 1; j <= u; j++ {
           var x int64
           fmt.Fscan(reader, &x)
           prefix[j] = prefix[j-1] + x
       }
       // skip rest
       for j := u + 1; j <= ti; j++ {
           var skip int64
           fmt.Fscan(reader, &skip)
       }
       // knapsack group: choose 0..u from this array
       // update dp in place reverse
       for j := k; j >= 1; j-- {
           // try take x from this array
           maxv := dp[j]
           // x from 1..u and x<=j
           lim := u
           if j < lim {
               lim = j
           }
           for x := 1; x <= lim; x++ {
               v := dp[j-x] + prefix[x]
               if v > maxv {
                   maxv = v
               }
           }
           dp[j] = maxv
       }
   }
   // result dp[k]
   fmt.Println(dp[k])
}
