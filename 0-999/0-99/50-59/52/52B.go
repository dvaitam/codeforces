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
   grid := make([]string, n)
   rowCount := make([]int, n)
   colCount := make([]int, m)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       grid[i] = s
       for j, ch := range s {
           if ch == '*' {
               rowCount[i]++
               colCount[j]++
           }
       }
   }
   var ans int64
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == '*' {
               rc := rowCount[i] - 1
               cc := colCount[j] - 1
               ans += int64(rc) * int64(cc)
           }
       }
   }
   fmt.Fprintln(writer, ans)
}
