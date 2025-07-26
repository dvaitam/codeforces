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
   for i := 0; i < t; i++ {
       var n, x, y int64
       fmt.Fscan(in, &n, &x, &y)
       // minimum overall place
       mn := x + y - n + 1
       if mn < 1 {
           mn = 1
       }
       if mn > n {
           mn = n
       }
       // maximum overall place
       mx := x + y - 1
       if mx > n {
           mx = n
       }
       fmt.Fprintln(out, mn, mx)
   }
}
