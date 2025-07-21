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

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var m, n int
   if _, err := fmt.Fscan(reader, &m, &n); err != nil {
       return
   }
   // dp[j]: next free time of painter j (0-indexed)
   dp := make([]int, n)
   // process each picture
   for i := 0; i < m; i++ {
       prev := 0
       for j := 0; j < n; j++ {
           var t int
           fmt.Fscan(reader, &t)
           start := max(dp[j], prev)
           finish := start + t
           dp[j] = finish
           prev = finish
       }
       // output finish time of last painter
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, prev)
   }
   // trailing newline
   writer.WriteByte('\n')
}
