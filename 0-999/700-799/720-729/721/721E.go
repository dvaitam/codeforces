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

   var L int64
   var n int
   var p, t int64
   _, err := fmt.Fscan(reader, &L, &n, &p, &t)
   if err != nil {
       return
   }
   var nextAvail int64 = 0
   var ans int64 = 0
   for i := 0; i < n; i++ {
       var li, ri int64
       fmt.Fscan(reader, &li, &ri)
       // earliest possible start in this segment
       y0 := li
       if nextAvail > y0 {
           y0 = nextAvail
       }
       // compute how many performances fit
       if ri - y0 >= p {
           // length available
           avail := ri - y0
           k := avail / p
           ans += k
           // update next available start after pause
           lastEnd := y0 + k*p
           nextAvail = lastEnd + t
       }
   }
   fmt.Fprint(writer, ans)
}
