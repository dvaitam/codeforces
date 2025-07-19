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

   var n, x, y int64
   if _, err := fmt.Fscan(in, &n, &x, &y); err != nil {
       return
   }
   u := y - n + 1
   if u <= 0 || u*u + n - 1 < x {
       fmt.Fprintln(out, -1)
       return
   }
   fmt.Fprintln(out, u)
   for i := int64(1); i < n; i++ {
       fmt.Fprintln(out, 1)
   }
}
