package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, d int
   fmt.Fscan(reader, &n, &m, &d)
   a := make([]int, m)
   b := make([]int, m)
   t := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &a[i], &b[i], &t[i])
       a[i]-- // zero-based
   }

   const INF64 = int64(4e18)
   dp := make([]int64, n)
   dp2 := make([]int64, n)
   // initial dp for first event
   for x := 0; x < n; x++ {
       dp[x] = int64(b[0] - abs(a[0]-x))
   }

   // temp deque
   deque := make([]int, n)
   // process further events
   for i := 1; i < m; i++ {
       dt := t[i] - t[i-1]
       k := dt * d
       // reset dp2
       if k >= n {
           // full range
           var best int64 = -INF64
           for x := 0; x < n; x++ {
               if dp[x] > best {
                   best = dp[x]
               }
           }
           for x := 0; x < n; x++ {
               dp2[x] = best + int64(b[i]-abs(a[i]-x))
           }
       } else {
           // sliding window max
           head, tail := 0, 0
           r := -1
           for x := 0; x < n; x++ {
               // expand right bound to x+k
               newR := min(n-1, x+k)
               for r < newR {
                   r++
                   // push r
                   for tail > head && dp[deque[tail-1]] <= dp[r] {
                       tail--
                   }
                   deque[tail] = r
                   tail++
               }
               // left bound
               l := max(0, x-k)
               // pop front if out of window
               for head < tail && deque[head] < l {
                   head++
               }
               // front is max
               best := dp[deque[head]]
               dp2[x] = best + int64(b[i]-abs(a[i]-x))
           }
       }
       // swap dp and dp2
       dp, dp2 = dp2, dp
   }

   // answer
   var ans int64 = -INF64
   for x := 0; x < n; x++ {
       if dp[x] > ans {
           ans = dp[x]
       }
   }
   fmt.Fprint(writer, ans)
}
