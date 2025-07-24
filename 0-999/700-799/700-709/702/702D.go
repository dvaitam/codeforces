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

   var d, k, a, b, t int64
   if _, err := fmt.Fscan(reader, &d, &k, &a, &b, &t); err != nil {
       return
   }

   var ans int64
   if d <= k {
       ans = d * a
   } else {
       ans = k * a
       rem := d - k
       full := rem / k
       // decide if repairing and driving full segments is beneficial
       if t + k*a < k*b {
           ans += full * (t + k*a)
       } else {
           // better to walk the remaining distance
           ans += rem * b
           fmt.Fprintln(writer, ans)
           return
       }
       // handle leftover distance
       rem2 := rem % k
       if rem2 > 0 {
           // either repair and drive rem2 km or walk rem2 km
           if t + rem2*a < rem2*b {
               ans += t + rem2*a
           } else {
               ans += rem2 * b
           }
       }
   }

   fmt.Fprintln(writer, ans)
}
