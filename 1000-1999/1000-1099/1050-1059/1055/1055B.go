package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, l int
   if _, err := fmt.Fscan(reader, &n, &m, &l); err != nil {
       return
   }
   h := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &h[i])
   }
   ctr := 0
   for i := 1; i <= n; i++ {
       if h[i] > l && h[i-1] <= l {
           ctr++
       }
   }
   for i := 0; i < m; i++ {
       var t int
       fmt.Fscan(reader, &t)
       if t == 0 {
           fmt.Fprintln(writer, ctr)
       } else {
           var p, q int
           fmt.Fscan(reader, &p, &q)
           if h[p] <= l && h[p]+q > l {
               left := h[p-1] > l
               right := h[p+1] > l
               if left && right {
                   ctr--
               } else if !left && !right {
                   ctr++
               }
           }
           h[p] += q
       }
   }
}
