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
       var n, m int
       fmt.Fscan(in, &n, &m)
       if n == 1 || m == 1 {
           fmt.Fprintln(out, "1 1")
       } else {
           fmt.Fprintln(out, "2 2")
       }
   }
}
