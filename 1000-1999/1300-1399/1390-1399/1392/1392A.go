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
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       same := true
       var first int
       for i := 0; i < n; i++ {
           var x int
           fmt.Fscan(reader, &x)
           if i == 0 {
               first = x
           } else if x != first {
               same = false
           }
       }
       if same {
           fmt.Fprintln(writer, n)
       } else {
           fmt.Fprintln(writer, 1)
       }
   }
}
