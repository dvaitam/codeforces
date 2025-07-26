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
       var x, y, a, b int64
       fmt.Fscan(reader, &x, &y, &a, &b)
       diff := y - x
       sum := a + b
       if diff%sum == 0 {
           fmt.Fprintln(writer, diff/sum)
       } else {
           fmt.Fprintln(writer, -1)
       }
   }
}
