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
       var row string
       fmt.Fscan(reader, &row)
       grid[i] = []byte(row)
   }
   ans := 0
   dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == 'W' {
               for _, d := range dirs {
                   ni, nj := i+d[0], j+d[1]
                   if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == 'P' {
                       ans++
                       grid[ni][nj] = '.'
                       break
                   }
               }
           }
       }
   }
   fmt.Println(ans)
}
