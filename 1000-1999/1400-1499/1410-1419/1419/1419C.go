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
   for tt := 0; tt < t; tt++ {
       var n, x int
       fmt.Fscan(in, &n, &x)
       sumDiff := 0
       allEq := true
       anyEq := false
       for i := 0; i < n; i++ {
           var a int
           fmt.Fscan(in, &a)
           if a != x {
               allEq = false
           } else {
               anyEq = true
           }
           sumDiff += a - x
       }
       switch {
       case allEq:
           fmt.Fprintln(out, 0)
       case anyEq || sumDiff == 0:
           fmt.Fprintln(out, 1)
       default:
           fmt.Fprintln(out, 2)
       }
   }
}
