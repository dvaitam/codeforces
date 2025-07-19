package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var l, r int
       var m int64
       fmt.Fscan(in, &l, &r, &m)
       for a := l; a <= r; a++ {
           // compute adjustment d
           d := (m - int64(l) + int64(r))%int64(a) + int64(l-r)
           if d > int64(r-l) {
               continue
           }
           var c int64
           if d > 0 {
               c = int64(l)
           } else {
               c = int64(l) - d
           }
           b := c + d
           fmt.Fprintf(out, "%d %d %d\n", a, b, c)
           break
       }
   }
}
