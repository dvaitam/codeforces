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
       var l1, r1, l2, r2 int
       fmt.Fscan(reader, &l1, &r1, &l2, &r2)
       x := l1
       y := l2
       if x == y {
           y++
       }
       fmt.Fprintln(writer, x, y)
   }
}
