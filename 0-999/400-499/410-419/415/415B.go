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
   var a, b int64
   if _, err := fmt.Fscan(reader, &n, &a, &b); err != nil {
       return
   }
   for i := 0; i < n; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       // maximum money M
       M := x * a / b
       // minimal tokens to return to get M dollars
       w := (M*b + a - 1) / a
       saved := x - w
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, saved)
   }
   fmt.Fprintln(writer)
}
