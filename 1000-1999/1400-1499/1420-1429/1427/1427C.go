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

   var r, n int
   fmt.Fscan(reader, &r, &n)
   t := make([]int, n+1)
   x := make([]int, n+1)
   y := make([]int, n+1)
   dp := make([]int, n+1)
   // initial position
   t[0] = 0
   x[0] = 1
   y[0] = 1
   // read events
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &t[i], &x[i], &y[i])
       dp[i] = -1_000_000_000
   }
   dp[0] = 0
   // window length: max distance between any two points in grid is 2*(r-1) <= 2*r
   window := 2 * r
   bestPref := -1_000_000_000
   L := 0
   ans := 0
   for i := 1; i <= n; i++ {
       // move L to maintain t[L] < t[i] - window
       for L < i && t[L] < t[i]-window {
           if dp[L] > bestPref {
               bestPref = dp[L]
           }
           L++
       }
       // use bestPref if exists
       if bestPref > -1_000_000_000 {
           dp[i] = bestPref + 1
       }
       // brute check remaining window
       for j := L; j < i; j++ {
           // if reachable
           if abs(x[i]-x[j])+abs(y[i]-y[j]) <= t[i]-t[j] {
               dp[i] = max(dp[i], dp[j]+1)
           }
       }
       if dp[i] > ans {
           ans = dp[i]
       }
   }
   fmt.Fprintln(writer, ans)
}
