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
   fmt.Fscan(reader, &n, &m)

   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       grid[i] = []byte(s)
   }

   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == '.' {
               if (i+j)%2 == 0 {
                   grid[i][j] = 'B'
               } else {
                   grid[i][j] = 'W'
               }
           }
       }
   }

   for i := 0; i < n; i++ {
       writer.Write(grid[i])
       writer.WriteByte('\n')
   }
}
