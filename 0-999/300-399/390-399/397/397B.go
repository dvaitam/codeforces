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
       var n, l, r int64
       fmt.Fscan(reader, &n, &l, &r)
       minX := (n + r - 1) / r
       maxX := n / l
       if minX <= maxX {
           fmt.Fprintln(writer, "Yes")
       } else {
           fmt.Fprintln(writer, "No")
       }
   }
}
