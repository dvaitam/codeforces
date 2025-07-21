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

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // total shifts per column
   total := make([]int, m)
   const INF = 1000000000
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       // check if row has any '1'
       has1 := false
       for _, ch := range s {
           if ch == '1' {
               has1 = true
               break
           }
       }
       if !has1 {
           fmt.Println(-1)
           return
       }
       // compute minimal shifts to bring '1' to each position
       dist := make([]int, m)
       for j := 0; j < m; j++ {
           dist[j] = INF
       }
       // forward scan for right shifts
       last := -INF
       for j := 0; j < 2*m; j++ {
           idx := j % m
           if s[idx] == '1' {
               last = j
           }
           if last >= 0 {
               d := j - last
               if d < dist[idx] {
                   dist[idx] = d
               }
           }
       }
       // backward scan for left shifts
       last = 2*m + INF
       for j := 2*m - 1; j >= 0; j-- {
           idx := j % m
           if s[idx] == '1' {
               last = j
           }
           d := last - j
           if d < dist[idx] {
               dist[idx] = d
           }
       }
       // accumulate
       for j := 0; j < m; j++ {
           total[j] += dist[j]
       }
   }
   // find minimal total shifts
   ans := INF
   for j := 0; j < m; j++ {
       if total[j] < ans {
           ans = total[j]
       }
   }
   if ans >= INF {
       fmt.Println(-1)
   } else {
       fmt.Println(ans)
   }
}
