package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = int64(4e18)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, s, e int
   fmt.Fscan(in, &n, &s, &e)
   x := make([]int64, n+1)
   a := make([]int64, n+1)
   b := make([]int64, n+1)
   c := make([]int64, n+1)
   d := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &x[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &b[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &c[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &d[i])
   }
   // ensure s < e by reflecting if needed
   if s > e {
       // mirror arrays
       ns, ne := n+1-s, n+1-e
       tx := make([]int64, n+1)
       ta := make([]int64, n+1)
       tb := make([]int64, n+1)
       tc := make([]int64, n+1)
       td := make([]int64, n+1)
       for i := 1; i <= n; i++ {
           j := n + 1 - i
           // mirror coordinate
           tx[i] = -x[j]
           ta[i] = b[j]
           tb[i] = a[j]
           tc[i] = d[j]
           td[i] = c[j]
       }
       x, a, b, c, d = tx, ta, tb, tc, td
       s, e = ns, ne
   }
   // dp[j]: min cost with j open segments
   dp := make([]int64, n+3)
   for i := range dp {
       dp[i] = INF
   }
   dp[0] = 0
   // process chairs
   for i := 1; i <= n; i++ {
       // add distance cost
       if i > 1 {
           dist := x[i] - x[i-1]
           for j := 0; j <= n; j++ {
               if dp[j] < INF {
                   dp[j] += dist * int64(j)
               }
           }
       }
       ndp := make([]int64, n+3)
       for j := range ndp {
           ndp[j] = INF
       }
       for j := 0; j <= n; j++ {
           v := dp[j]
           if v >= INF {
               continue
           }
           // at start
           if i == s {
               // open one outgoing
               ndp[j+1] = min(ndp[j+1], v + d[i])
           } else if i == e {
               // close one incoming
               if j >= 1 {
                   ndp[j-1] = min(ndp[j-1], v + a[i])
               }
           } else {
               // normal node: degree 2
               // open two outgoing
               if j+2 <= n {
                   ndp[j+2] = min(ndp[j+2], v + 2*d[i])
               }
               // one open one close
               if j >= 1 {
                   ndp[j] = min(ndp[j], v + d[i] + a[i])
               }
               // close two incoming
               if j >= 2 {
                   ndp[j-2] = min(ndp[j-2], v + 2*a[i])
               }
           }
       }
       dp = ndp
   }
   // result should be dp[0]
   fmt.Println(dp[0])
}
