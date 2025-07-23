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

   var A, B, n int64
   if _, err := fmt.Fscan(reader, &A, &B, &n); err != nil {
       return
   }
   for i := int64(0); i < n; i++ {
       var l, t, m int64
       fmt.Fscan(reader, &l, &t, &m)
       // height at l
       sl := A + (l-1)*B
       if sl > t {
           fmt.Fprintln(writer, -1)
           continue
       }
       // maximum r by height constraint: A + (r-1)*B <= t => r <= (t - A)/B + 1
       rMax := (t - A) / B + 1
       if rMax < l {
           fmt.Fprintln(writer, -1)
           continue
       }
       low, high := l, rMax
       var best int64 = l
       for low <= high {
           mid := (low + high) / 2
           k := mid - l + 1
           // s_mid
           sm := A + (mid-1)*B
           // sum = k*(sl + sm)/2
           sum := k * (sl + sm) / 2
           if sum <= m*t {
               best = mid
               low = mid + 1
           } else {
               high = mid - 1
           }
       }
       fmt.Fprintln(writer, best)
   }
}
