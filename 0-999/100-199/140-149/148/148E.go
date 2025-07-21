package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // dp[j]: max value taking j items so far
   const INF = int64(1e18)
   dp := make([]int64, m+1)
   for i := 1; i <= m; i++ {
       dp[i] = -INF
   }
   for i := 0; i < n; i++ {
       var k int
       fmt.Fscan(reader, &k)
       arr := make([]int64, k)
       for j := 0; j < k; j++ {
           fmt.Fscan(reader, &arr[j])
       }
       // prefix and suffix sums
       pre := make([]int64, k+1)
       suf := make([]int64, k+1)
       for j := 1; j <= k; j++ {
           pre[j] = pre[j-1] + arr[j-1]
           suf[j] = suf[j-1] + arr[k-j]
       }
       // best[j]: max sum taking j items from this shelf
       best := make([]int64, k+1)
       for j := 0; j <= k; j++ {
           var mx int64 = 0
           for l := 0; l <= j; l++ {
               v := pre[l] + suf[j-l]
               if v > mx {
                   mx = v
               }
           }
           best[j] = mx
       }
       // new dp
       newdp := make([]int64, m+1)
       for j := 0; j <= m; j++ {
           newdp[j] = -INF
       }
       for taken := 0; taken <= m; taken++ {
           if dp[taken] < 0 && taken > 0 {
               continue
           }
           // try take t items from this shelf
           maxT := k
           if remain := m - taken; remain < maxT {
               maxT = remain
           }
           for t := 0; t <= maxT; t++ {
               v := dp[taken] + best[t]
               if v > newdp[taken+t] {
                   newdp[taken+t] = v
               }
           }
       }
       dp = newdp
   }
   // answer
   res := dp[m]
   fmt.Println(res)
}
