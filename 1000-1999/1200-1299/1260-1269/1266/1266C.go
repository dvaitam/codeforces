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
   // special case 1x1
   if n == 1 && m == 1 {
       fmt.Fprintln(writer, 0)
       return
   }
   // if single row or single column
   if n == 1 || m == 1 {
       for i := 1; i <= n; i++ {
           for j := 1; j <= m; j++ {
               val := i + j
               if j == m {
                   fmt.Fprintln(writer, val)
               } else {
                   fmt.Fprint(writer, val, " ")
               }
           }
       }
       return
   }
   // general case
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           val := j * (i + m)
           if j == m {
               fmt.Fprintln(writer, val)
           } else {
               fmt.Fprint(writer, val, " ")
           }
       }
   }
}
