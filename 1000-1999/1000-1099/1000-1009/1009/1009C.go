package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   var sum int64
   for i := 0; i < m; i++ {
       var x, d int64
       fmt.Fscan(reader, &x, &d)
       sum += x * int64(n)
       if d >= 0 {
           // maximize distances: endpoints
           sum += d * int64(n-1) * int64(n) / 2
       } else {
           // minimize distances: choose median
           mid := n/2 + 1
           front := int64(mid - 1)
           last := int64(n - mid)
           sum += d * (front*(front+1)/2 + last*(last+1)/2)
       }
   }
   ans := float64(sum) / float64(n)
   writer := bufio.NewWriter(os.Stdout)
   fmt.Fprintf(writer, "%.15f", ans)
   writer.Flush()
}
