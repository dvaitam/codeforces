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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   const INF = int64(1) << 60
   minv, maxv := INF, -INF
   var cntMin, cntMax int64
   for i := 0; i < n; i++ {
       var b int64
       fmt.Fscan(reader, &b)
       if b < minv {
           minv = b
           cntMin = 1
       } else if b == minv {
           cntMin++
       }
       if b > maxv {
           maxv = b
           cntMax = 1
       } else if b == maxv {
           cntMax++
       }
   }
   diff := maxv - minv
   var ways int64
   if minv == maxv {
       ways = int64(n) * int64(n-1) / 2
   } else {
       ways = cntMin * cntMax
   }
   fmt.Fprintln(writer, diff, ways)
}
