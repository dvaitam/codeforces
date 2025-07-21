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
   thicks := make([]int, n)
   widths := make([]int, n)
   totalWidth := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &thicks[i], &widths[i])
       totalWidth += widths[i]
   }
   maxT := 2 * n
   // dp[t] = maximum sum of widths of vertical books with total thickness t
   const negInf = -1000000000
   dp := make([]int, maxT+1)
   for t := 1; t <= maxT; t++ {
       dp[t] = negInf
   }
   dp[0] = 0
   for i := 0; i < n; i++ {
       ti, wi := thicks[i], widths[i]
       for t := maxT; t >= ti; t-- {
           if prev := dp[t-ti]; prev >= 0 {
               if prev+wi > dp[t] {
                   dp[t] = prev + wi
               }
           }
       }
   }
   ans := maxT
   for t := 0; t <= maxT; t++ {
       if dp[t] >= 0 && dp[t]+t >= totalWidth {
           ans = t
           break
       }
   }
   fmt.Println(ans)
}
