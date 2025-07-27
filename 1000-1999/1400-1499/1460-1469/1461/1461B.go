package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       grid := make([]string, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &grid[i])
       }
       // Precompute consecutive stars to left and right for each cell
       left := make([][]int, n)
       right := make([][]int, n)
       for i := 0; i < n; i++ {
           left[i] = make([]int, m)
           right[i] = make([]int, m)
           // left
           for j := 0; j < m; j++ {
               if grid[i][j] == '*' {
                   if j > 0 {
                       left[i][j] = left[i][j-1] + 1
                   } else {
                       left[i][j] = 1
                   }
               }
           }
           // right
           for j := m - 1; j >= 0; j-- {
               if grid[i][j] == '*' {
                   if j < m-1 {
                       right[i][j] = right[i][j+1] + 1
                   } else {
                       right[i][j] = 1
                   }
               }
           }
       }
       var result int64 = 0
       // For each cell, count spruces with origin at (i,j)
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               if grid[i][j] != '*' {
                   continue
               }
               // try increasing height
               maxh := 0
               for k := 0; i+k < n; k++ {
                   // at row i+k, need at least k+1 stars to left and right
                   if left[i+k][j] >= k+1 && right[i+k][j] >= k+1 {
                       maxh++
                   } else {
                       break
                   }
               }
               result += int64(maxh)
           }
       }
       fmt.Fprintln(writer, result)
   }
}
