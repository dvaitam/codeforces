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
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var a, b int64
       fmt.Fscan(in, &a, &b)
       if b == 1 {
           fmt.Fprintln(out, "NO")
       } else {
           x := a * (b / 2)
           y := a * ((b + 1) / 2)
           if b%2 == 0 {
               y += a * b
           }
           fmt.Fprintln(out, "YES")
           fmt.Fprintf(out, "%d %d %d\n", x, y, x+y)
       }
   }
}
