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

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   fmt.Fscan(reader, &n, &k)
   var s string
   fmt.Fscan(reader, &s)
   var t string
   fmt.Fscan(reader, &t)
   t0 := t[0]
   t1 := t[1]
   // special case when both characters are same
   if t0 == t1 {
       cnt := 0
       for i := 0; i < n; i++ {
           if s[i] == t0 {
               cnt++
           }
       }
       // we can change up to k other chars to t0
       cnt = cnt + min(k, n-cnt)
       // number of pairs = cnt choose 2
       res := int64(cnt) * int64(cnt-1) / 2
       fmt.Fprint(writer, res)
       return
   }
   // dp[used_changes][count_t0] = max pairs
   const INF = -1_000_000_000
   // initialize dp
   dp := make([][]int, k+1)
   for i := 0; i <= k; i++ {
       dp[i] = make([]int, n+1)
       for j := 0; j <= n; j++ {
           dp[i][j] = INF
       }
   }
   dp[0][0] = 0
   // process positions
   for pos := 0; pos < n; pos++ {
       c := s[pos]
       // next dp
       ndp := make([][]int, k+1)
       for i := 0; i <= k; i++ {
           ndp[i] = make([]int, n+1)
           for j := 0; j <= n; j++ {
               ndp[i][j] = INF
           }
       }
       for used := 0; used <= k; used++ {
           for cnt0 := 0; cnt0 <= n; cnt0++ {
               cur := dp[used][cnt0]
               if cur < 0 {
                   continue
               }
               // keep original
               newUsed := used
               newCnt0 := cnt0
               add := 0
               if c == t0 {
                   newCnt0 = cnt0 + 1
               } else if c == t1 {
                   add = cnt0
               }
               ndp[newUsed][newCnt0] = max(ndp[newUsed][newCnt0], cur+add)
               // change to t0
               if used < k {
                   newUsed2 := used + 1
                   newCnt02 := cnt0 + 1
                   ndp[newUsed2][newCnt02] = max(ndp[newUsed2][newCnt02], cur)
                   // change to t1
                   newCnt1 := cnt0
                   add1 := cnt0
                   ndp[newUsed2][newCnt1] = max(ndp[newUsed2][newCnt1], cur+add1)
               }
           }
       }
       dp = ndp
   }
   // find best answer
   best := 0
   for used := 0; used <= k; used++ {
       for cnt0 := 0; cnt0 <= n; cnt0++ {
           if dp[used][cnt0] > best {
               best = dp[used][cnt0]
           }
       }
   }
   fmt.Fprint(writer, best)
}
