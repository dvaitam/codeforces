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

       if n >= 1 {
           fmt.Fprint(writer, 3)
       }
       if n >= 2 {
           fmt.Fprint(writer, " ", 5)
       }
       for idx := 2; idx < n; idx++ {
           fmt.Fprint(writer, " ", idx+4)
       }
       fmt.Fprintln(writer)
   }
}
