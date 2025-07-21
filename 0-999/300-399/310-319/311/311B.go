package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF = int64(4e18)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, p int
   fmt.Fscan(in, &n, &m, &p)
   pos := make([]int64, n+1)
   for i := 2; i <= n; i++ {
       var d int64
       fmt.Fscan(in, &d)
       pos[i] = pos[i-1] + d
   }
   ski := make([]int64, m+1)
   for i := 1; i <= m; i++ {
       var h int
       var t int64
       fmt.Fscan(in, &h, &t)
       ski[i] = t - pos[h]
   }
   sort.Slice(ski[1:], func(i, j int) bool { return ski[i+1] < ski[j+1] })
   // prefix sums of ski
   S := make([]int64, m+1)
   for i := 1; i <= m; i++ {
       S[i] = S[i-1] + ski[i]
   }
   dpPrev := make([]int64, m+1)
   dpCur := make([]int64, m+1)
   for i := 1; i <= m; i++ {
       dpPrev[i] = INF
   }
   dpPrev[0] = 0
   // cost function: cost for segment k+1..i
   cost := func(k, i int) int64 {
       // (i-k)*ski[i] - (S[i] - S[k])
       return int64(i-k)*ski[i] - (S[i] - S[k])
   }
   var compute func(int, int, int, int)
   compute = func(l, r, optl, optr int) {
       if l > r {
           return
       }
       mid := (l + r) >> 1
       dpCur[mid] = INF
       best := optl
       end := mid - 1
       if end > optr {
           end = optr
       }
       for k := optl; k <= end; k++ {
           v := dpPrev[k] + cost(k, mid)
           if v < dpCur[mid] {
               dpCur[mid] = v
               best = k
           }
       }
       // divide
       compute(l, mid-1, optl, best)
       compute(mid+1, r, best, optr)
   }
   for j := 1; j <= p; j++ {
       dpCur[0] = 0
       compute(1, m, 0, m)
       dpPrev, dpCur = dpCur, dpPrev
   }
   // answer is dpPrev[m]
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintln(out, dpPrev[m])
}
