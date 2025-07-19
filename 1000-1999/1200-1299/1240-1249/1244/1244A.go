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
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var a, b, c, d, k int
       fmt.Fscan(reader, &a, &b, &c, &d, &k)
       pen := (a + c - 1) / c
       pencil := (b + d - 1) / d
       if pen + pencil <= k {
           fmt.Fprintln(writer, pen, pencil)
       } else {
           fmt.Fprintln(writer, -1)
       }
   }
}
