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
   for t > 0 {
       t--
       var n int
       fmt.Fscan(in, &n)
       c0, c1, c2 := 0, 0, 0
       for i := 0; i < n; i++ {
           var a int
           fmt.Fscan(in, &a)
           switch a % 3 {
           case 0:
               c0++
           case 1:
               c1++
           case 2:
               c2++
           }
       }
       // Pair 1s and 2s to form zeros
       pairs := c1
       if c2 < pairs {
           pairs = c2
       }
       c0 += pairs
       c1 -= pairs
       c2 -= pairs
       // Combine triples of same remainder
       c0 += c1 / 3
       c0 += c2 / 3
       fmt.Fprintln(out, c0)
   }
}
