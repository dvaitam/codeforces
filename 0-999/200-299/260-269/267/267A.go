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

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   for i := 0; i < n; i++ {
       var a, b int64
       fmt.Fscan(in, &a, &b)
       var ops int64
       for a > 0 && b > 0 {
           if a >= b {
               ops += a / b
               a %= b
           } else {
               ops += b / a
               b %= a
           }
       }
       fmt.Fprintln(out, ops)
   }
}
