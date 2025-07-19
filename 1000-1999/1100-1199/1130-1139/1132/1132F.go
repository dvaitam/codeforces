package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)
   // remove consecutive duplicates
   seq := []byte{}
   for i := 0; i < len(s); i++ {
       b := s[i]
       if len(seq) == 0 || seq[len(seq)-1] != b {
           seq = append(seq, b)
       }
   }
   s = string(seq)
   n = len(seq)
   if n == 0 {
       fmt.Fprint(writer, 0)
       return
   }

   const inf = int(1e9)
   // dp[i][j]: min operations to clear s[i..j]
   dp := make([][]int, n)
   for i := 0; i < n; i++ {
       dp[i] = make([]int, n)
       for j := 0; j < n; j++ {
           dp[i][j] = inf
       }
       dp[i][i] = 1
   }
   for i := n - 1; i >= 0; i-- {
       for j := i + 1; j < n; j++ {
           // paint s[i] alone then rest
           dp[i][j] = dp[i+1][j] + 1
           // try merging with same character at k
           for k := i + 1; k <= j; k++ {
               if s[i] == s[k] {
                   cost := dp[k][j]
                   if k > i+1 {
                       cost += dp[i+1][k-1]
                   }
                   if cost < dp[i][j] {
                       dp[i][j] = cost
                   }
               }
           }
       }
   }
   fmt.Fprint(writer, dp[0][n-1])
}
