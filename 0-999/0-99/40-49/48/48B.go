package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([][]int, n)
   for i := 0; i < n; i++ {
       grid[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &grid[i][j])
       }
   }
   var a, b int
   fmt.Fscan(reader, &a, &b)

   // build prefix sums
   ps := make([][]int, n+1)
   for i := 0; i <= n; i++ {
       ps[i] = make([]int, m+1)
   }
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           ps[i][j] = grid[i-1][j-1] + ps[i-1][j] + ps[i][j-1] - ps[i-1][j-1]
       }
   }

   // function to compute sum in rectangle r..r+h-1, c..c+w-1
   rectSum := func(r, c, h, w int) int {
       // using 0-based grid, ps has offset
       r0, c0 := r, c
       r1, c1 := r+h, c+w
       return ps[r1][c1] - ps[r0][c1] - ps[r1][c0] + ps[r0][c0]
   }

   inf := n*m + 1
   best := inf
   // try orientation a x b
   if a <= n && b <= m {
       for i := 0; i + a <= n; i++ {
           for j := 0; j + b <= m; j++ {
               cnt := rectSum(i, j, a, b)
               if cnt < best {
                   best = cnt
               }
           }
       }
   }
   // try orientation b x a
   if b <= n && a <= m {
       for i := 0; i + b <= n; i++ {
           for j := 0; j + a <= m; j++ {
               cnt := rectSum(i, j, b, a)
               if cnt < best {
                   best = cnt
               }
           }
       }
   }
   if best == inf {
       best = 0
   }
   fmt.Println(best)
}
