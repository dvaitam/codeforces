package main

import (
   "bufio"
   "fmt"
   "os"
)

func maxInt16(a, b int16) int16 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s, p string
   fmt.Fscan(reader, &s)
   fmt.Fscan(reader, &p)
   n := len(s)
   m := len(p)
   // Precompute match positions
   match := make([]bool, n)
   for i := 0; i+m <= n; i++ {
       ok := true
       for j := 0; j < m; j++ {
           if s[i+j] != p[j] {
               ok = false
               break
           }
       }
       match[i] = ok
   }
   // dp[i][j]: max occurrences using first i chars with j removals
   // Use int16, initialize to -10000
   const negInf int16 = -10000
   dp := make([][]int16, n+1)
   for i := 0; i <= n; i++ {
       dp[i] = make([]int16, n+1)
       for j := 0; j <= n; j++ {
           dp[i][j] = negInf
       }
   }
   dp[0][0] = 0
   for i := 0; i <= n; i++ {
       for j := 0; j <= n; j++ {
           cur := dp[i][j]
           if cur < 0 {
               continue
           }
           if i < n {
               // remove s[i]
               if j+1 <= n {
                   dp[i+1][j+1] = maxInt16(dp[i+1][j+1], cur)
               }
               // keep s[i] without matching
               dp[i+1][j] = maxInt16(dp[i+1][j], cur)
           }
           // take match at i
           if i+m <= n && match[i] {
               dp[i+m][j] = maxInt16(dp[i+m][j], cur+1)
           }
       }
   }
   // Prepare answers
   ans := make([]int, n+1)
   // For each removals x, ans[x] = max dp[n][t] for t <= x
   best := int16(0)
   for x := 0; x <= n; x++ {
       if dp[n][x] > best {
           best = dp[n][x]
       }
       ans[x] = int(best)
   }
   // Output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprintf(writer, "%d", v)
   }
   writer.WriteByte('\n')
}
