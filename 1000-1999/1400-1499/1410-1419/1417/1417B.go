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
       var tt int64
       fmt.Fscan(reader, &n, &tt)
       var c int64
       for i := 0; i < n; i++ {
           var x int64
           fmt.Fscan(reader, &x)
           var r int
           if tt%2 == 0 && x == tt/2 {
               r = int(c % 2)
               c++
           } else if 2*x < tt {
               r = 0
           } else {
               r = 1
           }
           fmt.Fprint(writer, r, " ")
       }
       fmt.Fprintln(writer)
   }
}
