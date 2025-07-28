package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

// Solve the problem of finding, for each prefix of intervals [l_i, r_i],
// the maximum length of a non-decreasing subsequence one can achieve
// by choosing values a_i within the intervals.
func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   const INF int64 = 2e18
   const NEG_INF int64 = -2e18
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       // dp[len] = minimal possible tail value of a non-decreasing subsequence of length len
       dp := make([]int64, n+2)
       dp[0] = NEG_INF
       for i := 1; i <= n+1; i++ {
           dp[i] = INF
       }
       m := 0
       ans := make([]int, n)
       for i := 0; i < n; i++ {
           var l, r int64
           fmt.Fscan(reader, &l, &r)
           // find largest p such that dp[p] <= r
           lo, hi := 0, m+1
           for lo < hi {
               mid := (lo + hi) >> 1
               if dp[mid] <= r {
                   lo = mid + 1
               } else {
                   hi = mid
               }
           }
           p := lo - 1
           // new candidate tail for length p+1
           tail := dp[p]
           if l > tail {
               tail = l
           }
           if tail < dp[p+1] {
               dp[p+1] = tail
           }
           if p+1 > m {
               m = p + 1
           }
           ans[i] = m
       }
       // output answers for this test case
       for i, v := range ans {
           if i > 0 {
               writer.WriteByte(' ')
           }
           writer.WriteString(strconv.Itoa(v))
       }
       writer.WriteByte('\n')
   }
}
