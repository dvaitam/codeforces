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
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       if n >= 2 {
           for i := 2; i <= n; i++ {
               fmt.Fprint(writer, i)
               if i < n {
                   fmt.Fprint(writer, " ")
               } else {
                   fmt.Fprint(writer, " ")
               }
           }
       }
       fmt.Fprintln(writer, 1)
   }
}
