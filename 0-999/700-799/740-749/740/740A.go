package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, a, b, c int64
   if _, err := fmt.Fscan(reader, &n, &a, &b, &c); err != nil {
       return
   }
   // need k such that (n + k) % 4 == 0
   rem := (4 - n%4) % 4
   // dp[r] = min cost to get k mod 4 = r
   const INF = int64(4e18)
   dp := [4]int64{0, INF, INF, INF}
   // relax transitions enough times
   for it := 0; it < 4; it++ {
       old := dp
       for r := 0; r < 4; r++ {
           curr := old[r]
           if curr >= INF {
               continue
           }
           // buy one copybook
           nr := (r + 1) & 3
           dp[nr] = min(dp[nr], curr + a)
           // buy pack of two
           nr2 := (r + 2) & 3
           dp[nr2] = min(dp[nr2], curr + b)
           // buy pack of three
           nr3 := (r + 3) & 3
           dp[nr3] = min(dp[nr3], curr + c)
       }
   }
   // output minimal cost for remainder rem
   fmt.Println(dp[rem])
}
