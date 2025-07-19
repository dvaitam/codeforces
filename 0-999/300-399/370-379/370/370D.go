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

func pd(n, m, x, y, d int, e [][2]int, grid [][]byte) bool {
   // Check bounds
   if x < 0 || y < 0 || x+d >= n || y+d >= m {
       return false
   }
   // Check each point lies on border of the square
   for _, p := range e {
       xi, yi := p[0], p[1]
       if xi != x && xi != x+d && yi != y && yi != y+d {
           return false
       }
   }
   // Clone grid to mark
   tmp := make([][]byte, n)
   for i := 0; i < n; i++ {
       row := make([]byte, m)
       copy(row, grid[i])
       tmp[i] = row
   }
   // Mark borders
   for i := 0; i <= d; i++ {
       if tmp[x+i][y] == '.' {
           tmp[x+i][y] = '+'
       }
       if tmp[x][y+i] == '.' {
           tmp[x][y+i] = '+'
       }
       if tmp[x+i][y+d] == '.' {
           tmp[x+i][y+d] = '+'
       }
       if tmp[x+d][y+i] == '.' {
           tmp[x+d][y+i] = '+'
       }
   }
   // Output result
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i := 0; i < n; i++ {
       w.Write(tmp[i])
       w.WriteByte('\n')
   }
   return true
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   grid := make([][]byte, n)
   e := make([][2]int, 0)
   x0, y0 := n, m
   x1, y1 := -1, -1
   // Read grid and collect 'w'
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(in, &s)
       grid[i] = []byte(s)
       for j := 0; j < m; j++ {
           if grid[i][j] == 'w' {
               e = append(e, [2]int{i, j})
               if i < x0 {
                   x0 = i
               }
               if i > x1 {
                   x1 = i
               }
               if j < y0 {
                   y0 = j
               }
               if j > y1 {
                   y1 = j
               }
           }
       }
   }
   // If too many points
   if len(e) > n*4 {
       fmt.Println(-1)
       return
   }
   d := max(x1-x0, y1-y0)
   // Try possible squares
   if pd(n, m, x0, y0, d, e, grid) {
       return
   }
   if pd(n, m, x0, y1-d, d, e, grid) {
       return
   }
   if pd(n, m, x1-d, y0, d, e, grid) {
       return
   }
   fmt.Println(-1)
}
