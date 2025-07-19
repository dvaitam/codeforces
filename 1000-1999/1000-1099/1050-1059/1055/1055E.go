package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, s, m, k int
   fmt.Fscan(reader, &n, &s, &m, &k)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // segments and pr
   type segPair struct{ l, r int }
   seg := make([]segPair, s)
   pr := make([]int, n)
   for i := 0; i < n; i++ {
       pr[i] = i
   }
   for i := 0; i < s; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       l--
       r--
       seg[i] = segPair{l, r}
       start := l - 1
       if start < 0 {
           start = 0
       }
       for j := start; j <= r; j++ {
           if pr[j] < r {
               pr[j] = r
           }
       }
   }
   // prefix, b, dp inside check
   b := make([]int, n)
   pref := make([]int, n+1)

   // dp: use flat slice for speed
   dp := make([]int, n*(m+1))

   var check = func(val int) bool {
       // build b and pref
       for i := 0; i < n; i++ {
           if a[i] <= val {
               b[i] = 1
           } else {
               b[i] = 0
           }
       }
       pref[0] = 0
       for i := 1; i <= n; i++ {
           pref[i] = pref[i-1] + b[i-1]
       }
       // reset dp
       for i := range dp {
           dp[i] = 0
       }
       // initial segments
       for _, p := range seg {
           l, r := p.l, p.r
           // dp[r][1] = max(dp[r][1], sum b[l..r])
           idx := r*(m+1) + 1
           val0 := pref[r+1] - pref[l]
           if dp[idx] < val0 {
               dp[idx] = val0
           }
       }
       // dp transitions
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               cur := dp[i*(m+1)+j]
               if cur > 0 || j == 0 && dp[i*(m+1)+j] >= 0 {
                   // take segment starting at i
                   end := pr[i]
                   idx2 := end*(m+1) + (j + 1)
                   val2 := cur + (pref[end+1] - pref[i+1])
                   if dp[idx2] < val2 {
                       dp[idx2] = val2
                   }
               }
               // skip to next i
               if i+1 < n {
                   idx3 := (i+1)*(m+1) + j
                   if dp[idx3] < cur {
                       dp[idx3] = cur
                   }
               }
           }
       }
       // find best for j = m
       best := 0
       for i := 0; i < n; i++ {
           v := dp[i*(m+1)+m]
           if v > best {
               best = v
           }
       }
       return best >= k
   }

   // binary search
   lo, hi := 0, 1000000001
   for lo < hi {
       mid := lo + (hi-lo)/2
       if check(mid) {
           hi = mid
       } else {
           lo = mid + 1
       }
   }
   if lo == 1000000001 {
       lo = -1
   }
   fmt.Fprintln(writer, lo)
}
