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
   var n int
   var k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)
   // dp[i][l]: number of distinct subsequences of length l using prefix s[0:i]
   dp := make([][]int64, n+1)
   for i := range dp {
       dp[i] = make([]int64, n+1)
   }
   dp[0][0] = 1
   last := make([]int, 26)
   for i := 1; i <= n; i++ {
       c := s[i-1] - 'a'
       // copy counts for not taking s[i-1]
       for l := 0; l <= i; l++ {
           dp[i][l] = dp[i-1][l]
       }
       // take s[i-1]
       for l := 1; l <= i; l++ {
           added := dp[i-1][l-1]
           if last[c] > 0 {
               sub := dp[last[c]-1][l-1]
               added -= sub
           }
           if added < 0 {
               added = 0
           }
           dp[i][l] += added
           if dp[i][l] > k {
               dp[i][l] = k
           }
       }
       last[c] = i
   }
   // total distinct subsequences
   total := int64(0)
   for l := 0; l <= n; l++ {
       total += dp[n][l]
       if total > k {
           total = k + 1
       }
   }
   if total < k {
       fmt.Println(-1)
       return
   }
   // pick k subsequences with max lengths
   rem := k
   var ans int64
   for l := n; l >= 0 && rem > 0; l-- {
       cnt := dp[n][l]
       take := min(cnt, rem)
       ans += take * int64(n-l)
       rem -= take
   }
   fmt.Println(ans)
}
