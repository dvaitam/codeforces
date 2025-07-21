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
   var x int
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(in, &s)
       if len(s) >= 2 && s[1] == '+' {
           x++
       } else {
           x--
       }
   }
   fmt.Fprintln(out, x)
}
