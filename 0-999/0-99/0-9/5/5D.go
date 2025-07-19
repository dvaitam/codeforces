package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, v, l, d, w float64
   if _, err := fmt.Fscan(reader, &a, &v, &l, &d, &w); err != nil {
       return
   }
   t := 0.0
   if w > v {
       w = v
   }
   // distance to accelerate from 0 to w
   l3 := 0.5 * a * (w/a) * (w/a)
   var v1 float64
   if l3 <= d {
       // reach w before d
       t = w / a
       l1 := d - l3
       // compute possible max speed at end of segment
       l3 = math.Sqrt(a*l1 + w*w)
       if l3 > v {
           l3 = v
       }
       // accelerate from w to l3 and decelerate back to w symmetrically
       t += 2 * (l3 - w) / a
       // remaining distance at constant l3 speed
       t += (l1 - (l3*l3 - w*w) / a) / l3
       v1 = w
   } else {
       // cannot reach w before d
       t = math.Sqrt(2 * d / a)
       v1 = a * t
   }
   // remaining distance after d
   l -= d
   // distance to accelerate from v1 to v
   l2 := (v*v - v1*v1) / (2 * a)
   if l2 <= l {
       // can reach v
       t += (v - v1) / a
       t += (l - l2) / v
   } else {
       // cannot reach v
       t += (math.Sqrt(2*a*l + v1*v1) - v1) / a
   }
   fmt.Printf("%.6f\n", t)
}
