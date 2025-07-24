package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func reverse(s string) string {
   b := []byte(s)
   for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
       b[i], b[j] = b[j], b[i]
   }
   return string(b)
}

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   _, err := fmt.Fscan(reader, &n)
   if err != nil {
       return
   }
   costs := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &costs[i])
   }
   s := make([]string, n)
   sr := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &s[i])
       sr[i] = reverse(s[i])
   }
   const INF = int64(1e18)
   // dp[i][0]: min cost up to i, i not reversed
   // dp[i][1]: min cost up to i, i reversed
   var prev0, prev1 int64
   prev0 = 0
   prev1 = costs[0]
   for i := 1; i < n; i++ {
       cur0, cur1 := INF, INF
       // not reverse i
       if prev0 < INF && strings.Compare(s[i-1], s[i]) <= 0 {
           cur0 = min(cur0, prev0)
       }
       if prev1 < INF && strings.Compare(sr[i-1], s[i]) <= 0 {
           cur0 = min(cur0, prev1)
       }
       // reverse i
       if prev0 < INF && strings.Compare(s[i-1], sr[i]) <= 0 {
           cur1 = min(cur1, prev0+costs[i])
       }
       if prev1 < INF && strings.Compare(sr[i-1], sr[i]) <= 0 {
           cur1 = min(cur1, prev1+costs[i])
       }
       prev0, prev1 = cur0, cur1
   }
   res := prev0
   if prev1 < res {
       res = prev1
   }
   if res >= int64(1e18) {
       fmt.Println(-1)
   } else {
       fmt.Println(res)
   }
}
