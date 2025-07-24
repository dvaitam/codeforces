package main

import (
   "bufio"
   "fmt"
   "os"
)

func check(matrix [][]int, n, m, l int) bool {
   pref := make([][]int, n+1)
   for i := 0; i <= n; i++ {
      pref[i] = make([]int, m+1)
   }
   for i := 1; i <= n; i++ {
      for j := 1; j <= m; j++ {
         val := 0
         if matrix[i-1][j-1] >= l {
            val = 1
         }
         pref[i][j] = pref[i-1][j] + pref[i][j-1] - pref[i-1][j-1] + val
      }
   }
   target := l * l
   for i := l; i <= n; i++ {
      for j := l; j <= m; j++ {
         sum := pref[i][j] - pref[i-l][j] - pref[i][j-l] + pref[i-l][j-l]
         if sum == target {
            return true
         }
      }
   }
   return false
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   fmt.Fscan(in, &t)
   for ; t > 0; t-- {
      var n, m int
      fmt.Fscan(in, &n, &m)
      grid := make([][]int, n)
      for i := 0; i < n; i++ {
         row := make([]int, m)
         for j := 0; j < m; j++ {
            fmt.Fscan(in, &row[j])
         }
         grid[i] = row
      }
      low, high := 1, n
      if m < high {
         high = m
      }
      for low < high {
         mid := (low + high + 1) / 2
         if check(grid, n, m, mid) {
            low = mid
         } else {
            high = mid - 1
         }
      }
      fmt.Fprintln(out, low)
   }
}
