package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   var cs string
   if _, err := fmt.Fscan(reader, &n, &m, &cs); err != nil {
       return
   }
   pres := cs[0]
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   seen := make(map[byte]bool)
   dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] != pres {
               continue
           }
           for _, d := range dirs {
               ni, nj := i+d[0], j+d[1]
               if ni >= 0 && ni < n && nj >= 0 && nj < m {
                   c := grid[ni][nj]
                   if c != '.' && c != pres {
                       seen[c] = true
                   }
               }
           }
       }
   }
   fmt.Println(len(seen))
}
