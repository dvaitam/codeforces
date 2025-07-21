package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   p := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       p[i] = x
   }
   // prefix sums
   ps := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       ps[i] = ps[i-1] + p[i]
   }
   // segment sum of length m ending at j: seg[j] = ps[j] - ps[j-m]
   seg := make([]int64, n+1)
   for j := m; j <= n; j++ {
       seg[j] = ps[j] - ps[j-m]
   }
   // dp arrays: dpPrev for t-1, dpCurr for t
   const NEG_INF = -1000000000000000000
   dpPrev := make([]int64, n+1)
   dpCurr := make([]int64, n+1)
   // t = 1..k
   for t := 1; t <= k; t++ {
       // initialize dpCurr
       dpCurr[0] = NEG_INF
       for j := 1; j <= n; j++ {
           // skip j
           best := dpCurr[j-1]
           // take segment ending at j
           if j >= m {
               val := dpPrev[j-m] + seg[j]
               if val > best {
                   best = val
               }
           }
           dpCurr[j] = best
       }
       // swap dpPrev, dpCurr
       dpPrev, dpCurr = dpCurr, dpPrev
   }
   // result
   res := dpPrev[n]
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, res)
}
