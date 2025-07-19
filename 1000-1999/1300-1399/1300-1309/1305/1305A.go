package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
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
       a := make([]int, n)
       b := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
       }
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &b[i])
       }
       sort.Ints(a)
       sort.Ints(b)
       for i, v := range a {
           if i > 0 {
               fmt.Fprint(out, " ")
           }
           fmt.Fprint(out, v)
       }
       fmt.Fprintln(out)
       for i, v := range b {
           if i > 0 {
               fmt.Fprint(out, " ")
           }
           fmt.Fprint(out, v)
       }
       fmt.Fprintln(out)
   }
}
