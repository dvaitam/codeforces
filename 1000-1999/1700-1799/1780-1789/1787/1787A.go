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
       var n int
       fmt.Fscan(reader, &n)
       if n&1 == 1 {
           fmt.Fprintln(writer, -1)
       } else {
           fmt.Fprintf(writer, "%d %d\n", n/2, 1)
       }
   }
}
