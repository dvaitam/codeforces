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

   var n, m, q int
   fmt.Fscan(reader, &n, &m, &q)
   var s, t string
   fmt.Fscan(reader, &s, &t)

   crr := make([]int, n)
   // mark end positions where t matches s ending at i
   for i := m - 1; i < n; i++ {
       if s[i-m+1:i+1] == t {
           crr[i] = 1
       }
   }
   pre := make([]int, n)
   for i := 0; i < n; i++ {
       if i > 0 {
           pre[i] = pre[i-1] + crr[i]
       } else {
           pre[i] = crr[i]
       }
   }

   for qi := 0; qi < q; qi++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       start := l + m - 2
       end := r - 1
       if start > end || start < 0 || end < 0 {
           fmt.Fprintln(writer, 0)
           continue
       }
       if end >= n {
           end = n - 1
       }
       res := pre[end]
       if start > 0 {
           res -= pre[start-1]
       }
       fmt.Fprintln(writer, res)
   }
}
