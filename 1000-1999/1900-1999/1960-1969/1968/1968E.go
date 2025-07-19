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
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(in, &n)
       if n <= 3 {
           for i := 1; i < n; i++ {
               fmt.Fprintln(out, 1, i)
           }
           fmt.Fprintln(out, n, n)
       } else {
           fmt.Fprintln(out, 1, 1)
           for i := 3; i < n; i++ {
               fmt.Fprintln(out, 1, i)
           }
           fmt.Fprintln(out, n, n-1)
           fmt.Fprintln(out, n, n)
       }
       fmt.Fprintln(out)
   }
}
