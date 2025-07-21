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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var m int
   fmt.Fscan(reader, &m)
   var cur int64
   for i := 0; i < m; i++ {
       var w int
       var h int64
       fmt.Fscan(reader, &w, &h)
       if a[w] > cur {
           cur = a[w]
       }
       fmt.Fprintln(writer, cur)
       cur += h
   }
}
