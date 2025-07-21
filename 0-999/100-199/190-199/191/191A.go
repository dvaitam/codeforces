package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // dp[s][e]: max total length of path starting at s and ending at e
   const K = 26
   const INF = int64(1e18)
   dp := make([][K]int64, K)
   // initialize dp with very negative values
   for i := 0; i < K; i++ {
       for j := 0; j < K; j++ {
           dp[i][j] = -INF
       }
   }
   var ans int64 = 0
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       if len(s) == 0 {
           continue
       }
       u := int(s[0] - 'a')
       v := int(s[len(s)-1] - 'a')
       w := int64(len(s))
       // extend existing paths
       for st := 0; st < K; st++ {
           if dp[st][u] > -INF {
               if dp[st][u]+w > dp[st][v] {
                   dp[st][v] = dp[st][u] + w
               }
           }
       }
       // start new path of this single edge
       if w > dp[u][v] {
           dp[u][v] = w
       }
       // if this forms a cycle, update answer
       if dp[u][u] > ans {
           ans = dp[u][u]
       }
       if u != v {
           // also check dp[v][v] if path loops back via this
           if dp[v][v] > ans {
               ans = dp[v][v]
           }
       }
   }
   // print result
   if ans < 0 {
       ans = 0
   }
   fmt.Println(ans)
}
