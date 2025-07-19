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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       grid[i] = []byte(line)
   }

   // Directions: left, right, up, down
   dirs := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
   ok := true
   for i := 0; i < n && ok; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == 'W' {
               for _, d := range dirs {
                   ni, nj := i+d[0], j+d[1]
                   if ni >= 0 && ni < n && nj >= 0 && nj < m {
                       if grid[ni][nj] == 'S' {
                           ok = false
                           break
                       }
                       if grid[ni][nj] == '.' {
                           grid[ni][nj] = 'D'
                       }
                   }
               }
               if !ok {
                   break
               }
           }
       }
   }

   if !ok {
       fmt.Fprintln(writer, "No")
       return
   }
   fmt.Fprintln(writer, "Yes")
   for i := 0; i < n; i++ {
       writer.Write(grid[i])
       writer.WriteByte('\n')
   }
}
