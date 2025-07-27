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
       var l, r int64
       fmt.Fscan(in, &l, &r)
       if 2*l > r {
           fmt.Fprintln(out, "YES")
       } else {
           fmt.Fprintln(out, "NO")
       }
   }
}
