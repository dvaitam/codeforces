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
   for tt := 0; tt < t; tt++ {
       var n, q int
       fmt.Fscan(in, &n, &q)
       var s string
       fmt.Fscan(in, &s)
       // Process queries
       for i := 0; i < q; i++ {
           var l, r int
           fmt.Fscan(in, &l, &r)
           // check for same char as s[l-1] before l and same as s[r-1] after r
           first := s[l-1]
           last := s[r-1]
           ok := false
           for j := 0; j < l-1; j++ {
               if s[j] == first {
                   ok = true
                   break
               }
           }
           if !ok {
               for j := r; j < n; j++ {
                   if s[j] == last {
                       ok = true
                       break
                   }
               }
           }
           if ok {
               fmt.Fprintln(out, "YES")
           } else {
               fmt.Fprintln(out, "NO")
           }
       }
   }
