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
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       grid[i] = line
   }
   minR, maxR := n, -1
   minC, maxC := m, -1
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == '*' {
               if i < minR {
                   minR = i
               }
               if i > maxR {
                   maxR = i
               }
               if j < minC {
                   minC = j
               }
               if j > maxC {
                   maxC = j
               }
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := minR; i <= maxR; i++ {
       // substring from minC to maxC inclusive
       line := grid[i][minC : maxC+1]
       fmt.Fprintln(writer, line)
   }
}
