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
       var n int64
       fmt.Fscan(reader, &n)
       if n%2 == 1 {
           fmt.Fprintln(writer, "-1")
       } else {
           fmt.Fprintf(writer, "0 0 %d\n", n/2)
       }
   }
}
