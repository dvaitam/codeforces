package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k int
   fmt.Fscan(reader, &n, &m, &k)
   // read and ignore source positions
   for i := 0; i < k; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
   }
   // read and ignore destination positions
   for i := 0; i < k; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
   }

   var sb strings.Builder
   // move to top-left corner
   for i := 0; i < n-1; i++ {
       sb.WriteByte('U')
   }
   for i := 0; i < m-1; i++ {
       sb.WriteByte('L')
   }
   // snake traversal through all cells
   for i := 1; i <= n; i++ {
       if i&1 == 1 {
           for j := 0; j < m-1; j++ {
               sb.WriteByte('R')
           }
       } else {
           for j := 0; j < m-1; j++ {
               sb.WriteByte('L')
           }
       }
       // move down (stays if at bottom)
       sb.WriteByte('D')
   }

   ans := sb.String()
   fmt.Fprintln(writer, len(ans))
   fmt.Fprintln(writer, ans)
}
