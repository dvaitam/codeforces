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
   for ; t > 0; t-- {
       var a, b, c int64
       fmt.Fscan(reader, &a, &b, &c)
       sum := a + b + c
       if sum%9 != 0 {
           fmt.Fprintln(writer, "NO")
           continue
       }
       n := sum / 9
       if a < n || b < n || c < n {
           fmt.Fprintln(writer, "NO")
       } else {
           fmt.Fprintln(writer, "YES")
       }
   }
}
