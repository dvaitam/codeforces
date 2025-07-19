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

   var n, a, b int
   if _, err := fmt.Fscan(reader, &n, &a, &b); err != nil {
       return
   }
   x, y := -1, -1
   for i := 0; i <= n/a; i++ {
       rem := n - a*i
       if rem% b == 0 {
           x = i
           y = rem / b
           break
       }
   }
   if y == -1 {
       fmt.Fprint(writer, "-1")
       return
   }
   idx := 1
   for cnt := 0; cnt < x; cnt++ {
       l := idx
       r := idx + a - 1
       for j := l + 1; j <= r; j++ {
           fmt.Fprint(writer, j, " ")
       }
       fmt.Fprint(writer, l, " ")
       idx += a
   }
   for cnt := 0; cnt < y; cnt++ {
       l := idx
       r := idx + b - 1
       for j := l + 1; j <= r; j++ {
           fmt.Fprint(writer, j, " ")
       }
       fmt.Fprint(writer, l, " ")
       idx += b
   }
}
