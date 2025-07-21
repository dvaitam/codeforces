package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   // grid of painted pixels, 1-indexed with padding
   grid := make([][]bool, n+2)
   for i := range grid {
       grid[i] = make([]bool, m+2)
   }
   // process moves
   for t := 1; t <= k; t++ {
       var r, c int
       fmt.Fscan(in, &r, &c)
       grid[r][c] = true
       // check all 2x2 squares that include (r,c)
       for dr := -1; dr <= 0; dr++ {
           for dc := -1; dc <= 0; dc++ {
               x := r + dr
               y := c + dc
               if x >= 1 && x+1 <= n && y >= 1 && y+1 <= m {
                   if grid[x][y] && grid[x+1][y] && grid[x][y+1] && grid[x+1][y+1] {
                       fmt.Println(t)
                       return
                   }
               }
           }
       }
   }
   // no 2x2 black square formed
   fmt.Println(0)
}
