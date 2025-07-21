package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   ps := make([][]int64, n)
   for i := 0; i < n; i++ {
       ps[i] = make([]int64, m+1)
       for j := 1; j <= m; j++ {
           var v int64
           fmt.Fscan(reader, &v)
           ps[i][j] = ps[i][j-1] + v
       }
   }

   const INF = int64(-1e18)
   prev := make([]int64, m+2)
   curr := make([]int64, m+2)
   // row 1
   for j := 1; j <= m; j++ {
       prev[j] = ps[0][j]
   }

   // DP rows 2..n
   for i := 2; i <= n; i++ {
       // build prefix and suffix max of prev
       pref := make([]int64, m+2)
       suff := make([]int64, m+3)
       pref[0] = INF
       for j := 1; j <= m; j++ {
           pref[j] = max(pref[j-1], prev[j])
       }
       suff[m+1] = INF
       for j := m; j >= 1; j-- {
           suff[j] = max(suff[j+1], prev[j])
       }

       // compute curr
       row := ps[i-1]
       if i%2 == 0 {
           // prev index i-1 odd: prev c > curr j => suff[j+1]
           for j := 1; j <= m; j++ {
               val := suff[j+1]
               if val <= INF/2 {
                   curr[j] = INF
               } else {
                   curr[j] = row[j] + val
               }
           }
       } else {
           // i odd: prev i-1 even: prev c < curr j => pref[j-1]
           for j := 1; j <= m; j++ {
               val := pref[j-1]
               if val <= INF/2 {
                   curr[j] = INF
               } else {
                   curr[j] = row[j] + val
               }
           }
       }
       // swap prev and curr
       prev, curr = curr, prev
   }

   // answer is max over prev[1..m]
   ans := int64(-1e18)
   for j := 1; j <= m; j++ {
       if prev[j] > ans {
           ans = prev[j]
       }
   }
   fmt.Fprintln(writer, ans)
}
