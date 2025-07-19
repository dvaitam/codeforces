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
       var n int
       fmt.Fscan(in, &n)
       if n == 1 {
           fmt.Fprintln(out, 1)
           fmt.Fprintln(out, "1 2")
       } else if n == 2 {
           fmt.Fprintln(out, 1)
           fmt.Fprintln(out, "2 6")
       } else {
           m := n/2 + n%2
           fmt.Fprintln(out, m)
           for k := 0; k < m; k++ {
               i := 1 + 3*k
               j := 3*n - 3*k
               fmt.Fprintln(out, i, j)
           }
       }
   }
}
