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
       fmt.Fscan(in, &n)
       var s string
       fmt.Fscan(in, &s)
       res := make([]byte, n-1)
       a := int(s[0] - '0')
       for i := 1; i < n; i++ {
           b := int(s[i] - '0')
           if b == 1 {
               if a == 1 {
                   res[i-1] = '-'
                   a = 0
               } else {
                   res[i-1] = '+'
                   a = 1
               }
           } else {
               if a == 1 {
                   res[i-1] = '+'
                   // a remains 1
               } else {
                   res[i-1] = '-'
                   // a remains 0
               }
           }
       }
       fmt.Fprintln(out, string(res))
   }
}
