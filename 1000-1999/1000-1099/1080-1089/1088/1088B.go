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

   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   sort.Ints(a)
   s := 0
   idx := 0
   printed := 0
   for printed < k {
       if idx < n {
           d := a[idx] - s
           if d == 0 {
               idx++
               continue
           }
           fmt.Fprintln(out, d)
           s += d
           idx++
       } else {
           fmt.Fprintln(out, 0)
       }
       printed++
   }
}
