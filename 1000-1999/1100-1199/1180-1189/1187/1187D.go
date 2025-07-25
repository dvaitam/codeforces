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
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for ti := 0; ti < t; ti++ {
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
       ok := true
       for i := 0; i < n; i++ {
           if a[i] != b[i] {
               ok = false
               break
           }
       }
       if ok {
           fmt.Fprintln(out, "YES")
       } else {
           fmt.Fprintln(out, "NO")
       }
   }
}
