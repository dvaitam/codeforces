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
       var n int
       fmt.Fscan(in, &n)
       if n <= 2 {
           fmt.Fprintln(out, 0)
       } else {
           fmt.Fprintln(out, n-2)
       }
   }
}
