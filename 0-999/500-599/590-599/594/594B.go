package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   var dr, dv int
   if _, err := fmt.Fscan(in, &n, &dr, &dv); err != nil {
       return
   }
   r := float64(dr)
   v := float64(dv)
   for i := 0; i < n; i++ {
       var bb, ee int
       _, _ = fmt.Fscan(in, &bb, &ee)
       d := float64(ee - bb)
       // initial bounds for binary search on angle
       e := d / r
       b := e - 2*math.Pi
       if b < 0 {
           b = 0
       }
       // binary search for minimal angle satisfying path length >= d
       for j := 0; j < 50; j++ {
           mid := 0.5 * (b + e)
           t := math.Sin(mid * 0.5)
           if t < 0 {
               t = -t
           }
           t = r*mid + 2*r*t
           if t < d {
               b = mid
           } else {
               e = mid
           }
       }
       // time = (arc angle * radius) / velocity
       res := b * r / v
       fmt.Fprintf(out, "%.10f\n", res)
   }
}
