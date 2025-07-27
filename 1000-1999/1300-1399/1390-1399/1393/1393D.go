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
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       grid[i] = []byte(s)
   }
   // lft and rgt arrays
   lft := make([][]int, n)
   rgt := make([][]int, n)
   for i := 0; i < n; i++ {
       lft[i] = make([]int, m)
       rgt[i] = make([]int, m)
       // left to right
       for j := 0; j < m; j++ {
           if j > 0 && grid[i][j] == grid[i][j-1] {
               lft[i][j] = lft[i][j-1] + 1
           } else {
               lft[i][j] = 1
           }
       }
       // right to left
       for j := m - 1; j >= 0; j-- {
           if j+1 < m && grid[i][j] == grid[i][j+1] {
               rgt[i][j] = rgt[i][j+1] + 1
           } else {
               rgt[i][j] = 1
           }
       }
   }
   // d1: downward triangles (apex above, expanding downwards)
   d1 := make([][]int, n)
   for i := 0; i < n; i++ {
       d1[i] = make([]int, m)
       for j := 0; j < m; j++ {
           d1[i][j] = 1
           if i > 0 && grid[i][j] == grid[i-1][j] {
               // can extend previous triangle
               ext := min(d1[i-1][j], min(lft[i][j], rgt[i][j]) - 1)
               if ext >= 0 {
                   d1[i][j] = ext + 1
               }
           }
       }
   }
   // d2: upward triangles (apex below, expanding upwards)
   d2 := make([][]int, n)
   for i := n - 1; i >= 0; i-- {
       d2[i] = make([]int, m)
       for j := 0; j < m; j++ {
           d2[i][j] = 1
           if i+1 < n && grid[i][j] == grid[i+1][j] {
               ext := min(d2[i+1][j], min(lft[i][j], rgt[i][j]) - 1)
               if ext >= 0 {
                   d2[i][j] = ext + 1
               }
           }
       }
   }
   var ans int64
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           // number of diamonds centered at (i,j)
           cnt := d1[i][j]
           if d2[i][j] < cnt {
               cnt = d2[i][j]
           }
           ans += int64(cnt)
       }
   }
   fmt.Fprint(writer, ans)
}
