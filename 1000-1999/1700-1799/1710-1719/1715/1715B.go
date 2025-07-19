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

   var t int
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var n int
       var k, b, s int64
       fmt.Fscan(reader, &n, &k, &b, &s)
       base := b * k
       minSum := base
       maxSum := base + int64(n)*(k-1)
       if s < minSum || s > maxSum {
           fmt.Fprintln(writer, -1)
           continue
       }
       rem := s - base
       a := make([]int64, n)
       a[0] = base
       for i := 1; i < n && rem > 0; i++ {
           add := k - 1
           if rem < add {
               add = rem
           }
           a[i] = add
           rem -= add
       }
       a[0] += rem
       for i := 0; i < n; i++ {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, a[i])
       }
       writer.WriteByte('\n')
   }
}
