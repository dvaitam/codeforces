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
   for tc := 0; tc < t; tc++ {
       var n int
       var s string
       fmt.Fscan(in, &n, &s)
       cnt := 0
       run := 0
       for i := 0; i < n; i++ {
           if s[i] == 'x' {
               run++
               if run >= 3 {
                   cnt++
               }
           } else {
               run = 0
           }
       }
       fmt.Fprintln(out, cnt)
   }
}
