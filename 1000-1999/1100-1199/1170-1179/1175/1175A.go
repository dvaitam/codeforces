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
       var n, k uint64
       fmt.Fscan(in, &n, &k)
       var steps uint64
       if k <= 1 {
           steps = n
       } else {
           for n > 0 {
               if n < k {
                   steps += n
                   break
               }
               r := n % k
               if r != 0 {
                   steps += r
                   n -= r
               } else {
                   n /= k
                   steps++
               }
           }
       }
       fmt.Fprintln(out, steps)
   }
}
