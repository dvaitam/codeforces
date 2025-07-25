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
   var L int64
   fmt.Fscan(reader, &n, &L)
   per := L / int64(n)
   res := make([][2]int64, n)
   const INF = int64(1000000000000000000)
   var prevR int64 = 0
   for i := 0; i < n; i++ {
       li := prevR
       // query f_i(li)
       var fl int64
       if li == 0 {
           fl = 0
       } else {
           fmt.Fprintf(writer, "? %d %d\n", i+1, li)
           writer.Flush()
           fmt.Fscan(reader, &fl)
       }
       target := fl + per
       // binary search r in [li, INF]
       lo, hi := li, INF
       for lo < hi {
           mid := lo + (hi-lo)/2
           fmt.Fprintf(writer, "? %d %d\n", i+1, mid)
           writer.Flush()
           var fm int64
           fmt.Fscan(reader, &fm)
           if fm >= target {
               hi = mid
           } else {
               lo = mid + 1
           }
       }
       ri := lo
       res[i][0], res[i][1] = li, ri
       prevR = ri
   }
   // output result
   fmt.Fprintln(writer, "!")
   for i := 0; i < n; i++ {
       fmt.Fprintf(writer, "%d %d\n", res[i][0], res[i][1])
   }
}
