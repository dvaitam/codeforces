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

   var q int
   if _, err := fmt.Fscan(reader, &q); err != nil {
       return
   }
   for i := 0; i < q; i++ {
       var a, b, c int64
       fmt.Fscan(reader, &a, &b, &c)
       sum := a + b + c
       fmt.Fprintln(writer, sum/2)
   }
}
