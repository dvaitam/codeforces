package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       grid[i] = []byte(s)
   }
   // Precompute up, down, right arrays
   up := make([][]int, n)
   down := make([][]int, n)
   right := make([][]int, n)
   for i := 0; i < n; i++ {
       up[i] = make([]int, m)
       down[i] = make([]int, m)
       right[i] = make([]int, m)
   }
   // up
   for j := 0; j < m; j++ {
       for i := 0; i < n; i++ {
           if i > 0 && grid[i][j] == grid[i-1][j] {
               up[i][j] = up[i-1][j] + 1
           } else {
               up[i][j] = 1
           }
       }
   }
   // down
   for j := 0; j < m; j++ {
       for i := n - 1; i >= 0; i-- {
           if i+1 < n && grid[i][j] == grid[i+1][j] {
               down[i][j] = down[i+1][j] + 1
           } else {
               down[i][j] = 1
           }
       }
   }
   // right
   for i := 0; i < n; i++ {
       for j := m - 1; j >= 0; j-- {
           if j+1 < m && grid[i][j] == grid[i][j+1] {
               right[i][j] = right[i][j+1] + 1
           } else {
               right[i][j] = 1
           }
       }
   }
   var ans int64
   // For each column as left boundary
   for j := 0; j < m; j++ {
       for i := 0; i < n; i++ {
           // middle stripe ends at i
           h2 := up[i][j]
           midEnd := i
           topEnd := midEnd - h2
           if topEnd < 0 {
               continue
           }
           botStart := midEnd + 1
           if botStart >= n {
               continue
           }
           h1 := up[topEnd][j]
           h3 := down[botStart][j]
           c2 := grid[midEnd][j]
           c1 := grid[topEnd][j]
           c3 := grid[botStart][j]
           if c1 == c2 || c2 == c3 {
               continue
           }
           // stripe height
           h := h2
           if h1 < h {
               h = h1
           }
           if h3 < h {
               h = h3
           }
           // compute minimal width over stripes
           minW := m // maximum possible
           for u := 0; u < h; u++ {
               // rows: top: topEnd-u, mid: midEnd-u, bot: botStart+u
               rTop := topEnd - u
               if right[rTop][j] < minW {
                   minW = right[rTop][j]
               }
               rMid := midEnd - u
               if right[rMid][j] < minW {
                   minW = right[rMid][j]
               }
               rBot := botStart + u
               if right[rBot][j] < minW {
                   minW = right[rBot][j]
               }
           }
           ans += int64(minW)
       }
   }
   fmt.Println(ans)
}
