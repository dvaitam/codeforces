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
   prefRow := make([][]int, n)
   for i := 0; i < n; i++ {
      prefRow[i] = make([]int, m)
      sum := 0
      for j := 0; j < m; j++ {
         if grid[i][j] == 1 {
            sum++
         }
         prefRow[i][j] = sum
      }
   }
   prefCol := make([][]int, n)
   for i := 0; i < n; i++ {
      prefCol[i] = make([]int, m)
   }
   for j := 0; j < m; j++ {
      sum := 0
      for i := 0; i < n; i++ {
         if grid[i][j] == 1 {
            sum++
         }
         prefCol[i][j] = sum
      }
   }
   ans := 0
   for i := 0; i < n; i++ {
      for j := 0; j < m; j++ {
         if grid[i][j] == 1 {
            continue
         }
         if j > 0 && prefRow[i][j-1] > 0 {
            ans++
         }
         if j+1 < m && prefRow[i][m-1]-prefRow[i][j] > 0 {
            ans++
         }
         if i > 0 && prefCol[i-1][j] > 0 {
            ans++
         }
         if i+1 < n && prefCol[n-1][j]-prefCol[i][j] > 0 {
            ans++
         }
      }
   }
   fmt.Println(ans)
}
