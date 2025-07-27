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

   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n, p, k int
       var x, y int64
       fmt.Fscan(in, &n, &p, &k)
       var s string
       fmt.Fscan(in, &s)
       fmt.Fscan(in, &x, &y)
       // 1-indexed array of bytes
       a := make([]byte, n+1)
       for i := 1; i <= n; i++ {
           a[i] = s[i-1]
       }
       // dp[i] = number of zeros at positions i, i+k, i+2k, ... <=n
       dp := make([]int64, n+2)
       for i := n; i >= 1; i-- {
           next := i + k
           var cnt int64
           if a[i] == '0' {
               cnt = 1
           }
           if next <= n {
               dp[i] = cnt + dp[next]
           } else {
               dp[i] = cnt
           }
       }
       // consider deletions d from 0 to n-p
       var ans int64 = 1<<62
       // we need final p-th cell exists => d <= n-p
       for d := 0; d <= n-p; d++ {
           idx := p + d
           // dp at idx
           cost := int64(d)*y + dp[idx]*x
           if cost < ans {
               ans = cost
           }
       }
       fmt.Fprintln(out, ans)
   }
}
