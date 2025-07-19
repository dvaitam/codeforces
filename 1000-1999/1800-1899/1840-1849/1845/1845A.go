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
   for tc := 0; tc < t; tc++ {
       var n, k, x int
       fmt.Fscan(in, &n, &k, &x)
       if x != 1 {
           fmt.Fprintln(out, "YES")
           fmt.Fprintln(out, n)
           for i := 0; i < n; i++ {
               if i > 0 {
                   out.WriteString(" ")
               }
               out.WriteString("1")
           }
           fmt.Fprintln(out)
       } else {
           if k == 1 {
               fmt.Fprintln(out, "NO")
           } else if n%2 == 0 {
               fmt.Fprintln(out, "YES")
               fmt.Fprintln(out, n/2)
               for i := 0; i < n/2; i++ {
                   if i > 0 {
                       out.WriteString(" ")
                   }
                   out.WriteString("2")
               }
               fmt.Fprintln(out)
           } else if k == 2 {
               fmt.Fprintln(out, "NO")
           } else if n == 1 {
               fmt.Fprintln(out, "NO")
           } else {
               // odd n, k > 2, n > 1
               fmt.Fprintln(out, "YES")
               fmt.Fprintln(out, n/2)
               // print n/2 - 1 times 2, then 3
               for i := 0; i < n/2-1; i++ {
                   out.WriteString("2 ")
               }
               out.WriteString("3")
               fmt.Fprintln(out)
           }
       }
   }
}
