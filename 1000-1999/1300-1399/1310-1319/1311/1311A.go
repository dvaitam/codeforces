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
   for i := 0; i < t; i++ {
       var a, b int64
       fmt.Fscan(reader, &a, &b)
       switch {
       case a == b:
           fmt.Fprintln(writer, 0)
       case a < b:
           d := b - a
           if d%2 == 1 {
               fmt.Fprintln(writer, 1)
           } else {
               fmt.Fprintln(writer, 2)
           }
       default:
           d := a - b
           if d%2 == 0 {
               fmt.Fprintln(writer, 1)
           } else {
               fmt.Fprintln(writer, 2)
           }
       }
   }
}
