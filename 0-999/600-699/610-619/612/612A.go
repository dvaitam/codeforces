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

   var n, p, q int
   if _, err := fmt.Fscan(reader, &n, &p, &q); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)

   xx, yy := -1, -1
   for x := 0; x*q <= n; x++ {
       rem := n - x*q
       if rem%p == 0 {
           xx = x
           yy = rem / p
           break
       }
   }
   if xx == -1 && yy == -1 {
       fmt.Fprintln(writer, -1)
       return
   }

   total := xx + yy
   fmt.Fprintln(writer, total)
   idx := 0
   for i := 0; i < xx; i++ {
       end := idx + q
       fmt.Fprintln(writer, s[idx:end])
       idx = end
   }
   for i := 0; i < yy; i++ {
       end := idx + p
       fmt.Fprintln(writer, s[idx:end])
       idx = end
   }
}
