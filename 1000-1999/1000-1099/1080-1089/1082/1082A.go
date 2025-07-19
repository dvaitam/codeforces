package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var n, x, y, d int64
       fmt.Fscan(in, &n, &x, &y, &d)
       diff := abs(y - x)
       if diff % d == 0 {
           fmt.Fprintln(out, diff / d)
       } else {
           const inf = int64(1) << 62
           ans := inf
           // try via 1
           if (y-1) % d == 0 {
               steps := (y-1)/d + (x-1 + d - 1)/d
               if steps < ans {
                   ans = steps
               }
           }
           // try via n
           if (n - y) % d == 0 {
               steps := (n-y)/d + (n-x + d - 1)/d
               if steps < ans {
                   ans = steps
               }
           }
           if ans == inf {
               fmt.Fprintln(out, -1)
           } else {
               fmt.Fprintln(out, ans)
           }
       }
   }
}
