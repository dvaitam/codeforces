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

   var n int
   fmt.Fscan(reader, &n)
   grid := make([]string, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   // prefix sum of black cells
   sum := make([][]int, n+1)
   for i := range sum {
       sum[i] = make([]int, n+1)
   }
   for i := 1; i <= n; i++ {
       for j := 1; j <= n; j++ {
           add := 0
           if grid[i][j-1] == '#' {
               add = 1
           }
           sum[i][j] = sum[i-1][j] + sum[i][j-1] - sum[i-1][j-1] + add
       }
   }
   // dp[x1][x2][y1][y2]
   var dp [51][51][51][51]int
   // iterate over sizes
   for h := 1; h <= n; h++ {
       for w := 1; w <= n; w++ {
           for x1 := 1; x1+h-1 <= n; x1++ {
               x2 := x1 + h - 1
               for y1 := 1; y1+w-1 <= n; y1++ {
                   y2 := y1 + w - 1
                   // check if any black in [x1..x2][y1..y2]
                   cnt := sum[x2][y2] - sum[x1-1][y2] - sum[x2][y1-1] + sum[x1-1][y1-1]
                   if cnt == 0 {
                       dp[x1][x2][y1][y2] = 0
                   } else {
                       // initial cost: paint whole rectangle
                       best := max(h, w)
                       // horizontal splits
                       for k := x1; k < x2; k++ {
                           cost := dp[x1][k][y1][y2] + dp[k+1][x2][y1][y2]
                           best = min(best, cost)
                       }
                       // vertical splits
                       for k := y1; k < y2; k++ {
                           cost := dp[x1][x2][y1][k] + dp[x1][x2][k+1][y2]
                           best = min(best, cost)
                       }
                       dp[x1][x2][y1][y2] = best
                   }
               }
           }
       }
   }
   res := dp[1][n][1][n]
   fmt.Fprintln(writer, res)
}
