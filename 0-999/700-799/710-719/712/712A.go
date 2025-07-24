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
   fmt.Fscan(reader, &n)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   b := make([]int64, n)
   for i := 0; i < n-1; i++ {
       b[i] = a[i] + a[i+1]
   }
   b[n-1] = a[n-1]
   for i := 0; i < n; i++ {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, b[i])
   }
   fmt.Fprintln(writer)
}
