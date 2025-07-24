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
   fmt.Fscan(reader, &n)
   const INF = int64(1e18)
   low := -INF
   high := INF
   delta := int64(0)
   for i := 0; i < n; i++ {
       var c int64
       var d int
       fmt.Fscan(reader, &c, &d)
       if d == 1 {
           // division 1: rating before contest >= 1900
           if v := 1900 - delta; int64(v) > low {
               low = int64(v)
           }
       } else {
           // division 2: rating before contest <= 1899
           if v := 1899 - delta; int64(v) < high {
               high = int64(v)
           }
       }
       delta += c
   }
   if low > high {
       fmt.Fprintln(writer, "Impossible")
   } else if high == INF {
       fmt.Fprintln(writer, "Infinity")
   } else {
       fmt.Fprintln(writer, high + delta)
   }
}
