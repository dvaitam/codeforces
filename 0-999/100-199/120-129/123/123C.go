package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   var k int64
   fmt.Fscan(in, &n, &m, &k)
   D := n + m - 1
   // read priorities and build list
   type cell struct{ p int; d int }
   cells := make([]cell, 0, n*m)
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           var p int
           fmt.Fscan(in, &p)
           cells = append(cells, cell{p: p, d: i + j})
       }
   }
   sort.Slice(cells, func(i, j int) bool { return cells[i].p < cells[j].p })
   seen := make([]bool, D)
   order := make([]int, 0, D)
   for _, c := range cells {
       if !seen[c.d] {
           seen[c.d] = true
           order = append(order, c.d)
       }
   }
   // fixed positions: 0=free, 1='(', 2=')'
   fixed := make([]int, D)
   // DP arrays
   // Cap counts at k
   var countWays func() int64
   countWays = func() int64 {
       // dp[balance]
       dp := make([]int64, D+2)
       dp[0] = 1
       for pos := 0; pos < D; pos++ {
           ndp := make([]int64, D+2)
           if fixed[pos] == 1 {
               // must '('
               for bal := 0; bal <= D; bal++ {
                   v := dp[bal]
                   if v == 0 {
                       continue
                   }
                   nb := bal + 1
                   ndp[nb] += v
                   if ndp[nb] > k {
                       ndp[nb] = k
                   }
               }
           } else if fixed[pos] == 2 {
               // must ')'
               for bal := 1; bal <= D; bal++ {
                   v := dp[bal]
                   if v == 0 {
                       continue
                   }
                   nb := bal - 1
                   ndp[nb] += v
                   if ndp[nb] > k {
                       ndp[nb] = k
                   }
               }
           } else {
               // free: try both
               for bal := 0; bal <= D; bal++ {
                   v := dp[bal]
                   if v == 0 {
                       continue
                   }
                   // '('
                   nb := bal + 1
                   ndp[nb] += v
                   if ndp[nb] > k {
                       ndp[nb] = k
                   }
                   // ')'
                   if bal > 0 {
                       nb2 := bal - 1
                       ndp[nb2] += v
                       if ndp[nb2] > k {
                           ndp[nb2] = k
                       }
                   }
               }
           }
           dp = ndp
       }
       return dp[0]
   }
   // Greedy assign
   for _, d := range order {
       // try '('
       fixed[d] = 1
       cnt := countWays()
       if cnt < k {
           k -= cnt
           fixed[d] = 2
       }
   }
   // output grid
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       row := make([]byte, m)
       for j := 0; j < m; j++ {
           if fixed[i+j] == 1 {
               row[j] = '('
           } else {
               row[j] = ')'
           }
       }
       grid[i] = row
   }
   for i := 0; i < n; i++ {
       out.Write(grid[i])
       out.WriteByte('\n')
   }
}
