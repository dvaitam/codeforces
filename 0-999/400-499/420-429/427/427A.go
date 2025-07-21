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
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var available, untreated int
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x == -1 {
           if available > 0 {
               available--
           } else {
               untreated++
           }
       } else {
           available += x
       }
   }
   fmt.Fprintln(writer, untreated)
}
