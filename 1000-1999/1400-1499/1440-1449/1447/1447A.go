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
       fmt.Fprintln(writer, n)
       for j := 1; j <= n; j++ {
           fmt.Fprint(writer, j, " ")
       }
       fmt.Fprintln(writer)
   }
}
