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
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       grid[i] = []byte(s)
   }
   // Count occurrences in rows and columns
   rowCount := make([][26]int, n)
   colCount := make([][26]int, m)
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           c := grid[i][j] - 'a'
           rowCount[i][c]++
           colCount[j][c]++
       }
   }
   // Collect non-repeated letters
   var ans []byte
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           c := grid[i][j] - 'a'
           if rowCount[i][c] == 1 && colCount[j][c] == 1 {
               ans = append(ans, grid[i][j])
           }
       }
   }
   // Output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   if len(ans) > 0 {
       writer.Write(ans)
   }
   writer.WriteByte('\n')
}
