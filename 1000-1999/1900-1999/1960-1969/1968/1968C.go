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
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n-1)
       for i := 0; i < n-1; i++ {
           fmt.Fscan(reader, &a[i])
       }
       b := make([]int, n)
       b[0] = 501
       for i := 1; i < n; i++ {
           b[i] = b[i-1] + a[i-1]
       }
       for i := 0; i < n; i++ {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprintf(writer, "%d", b[i])
       }
       writer.WriteByte('\n')
   }
}
