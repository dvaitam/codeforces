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
   // grid of bytes, 0 means empty
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       grid[i] = make([]byte, m)
   }
   // process in row-major order
   for r := 0; r < n; r++ {
       for c := 0; c < m; c++ {
           if grid[r][c] != 0 {
               continue
           }
           // find smallest color
           var color byte
           for ch := byte('A'); ch <= 'Z'; ch++ {
               // immediate neighbors: up and left
               ok := true
               if r > 0 && grid[r-1][c] == ch {
                   ok = false
               }
               if c > 0 && grid[r][c-1] == ch {
                   ok = false
               }
               if ok {
                   color = ch
                   break
               }
           }
           // find max size
           maxs := 0
           // limit by boundaries: max square size fitting
           limit := n - r
           if m-c < limit {
               limit = m - c
           }
           for s := 1; s <= limit; s++ {
               bad := false
               // check inside emptiness
               for i := r; i < r+s && !bad; i++ {
                   for j := c; j < c+s; j++ {
                       if grid[i][j] != 0 {
                           bad = true
                           break
                       }
                   }
               }
               if bad {
                   break
               }
               // check adjacency borders for same color
               // top border
               if r-1 >= 0 {
                   for j := c; j < c+s; j++ {
                       if grid[r-1][j] == color {
                           bad = true
                           break
                       }
                   }
               }
               if bad {
                   break
               }
               // left border
               if c-1 >= 0 {
                   for i := r; i < r+s; i++ {
                       if grid[i][c-1] == color {
                           bad = true
                           break
                       }
                   }
               }
               if bad {
                   break
               }
               // bottom border
               if r+s < n {
                   for j := c; j < c+s; j++ {
                       if grid[r+s][j] == color {
                           bad = true
                           break
                       }
                   }
               }
               if bad {
                   break
               }
               // right border
               if c+s < m {
                   for i := r; i < r+s; i++ {
                       if grid[i][c+s] == color {
                           bad = true
                           break
                       }
                   }
               }
               if bad {
                   break
               }
               maxs = s
           }
           // fill s = maxs
           for i := r; i < r+maxs; i++ {
               for j := c; j < c+maxs; j++ {
                   grid[i][j] = color
               }
           }
       }
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < n; i++ {
       writer.Write(grid[i])
       writer.WriteByte('\n')
   }
}
