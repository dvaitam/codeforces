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

   var t int
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       grid := make([]string, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &grid[i])
       }
       ans := 0
       // Cells in last row (row n-1), columns 0..m-2 should be 'R'
       for j := 0; j < m-1; j++ {
           if grid[n-1][j] != 'R' {
               ans++
           }
       }
       // Cells in last column (col m-1), rows 0..n-2 should be 'D'
       for i := 0; i < n-1; i++ {
           if grid[i][m-1] != 'D' {
               ans++
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
