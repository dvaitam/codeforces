package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, a, d int
   if _, err := fmt.Fscan(reader, &n, &a, &d); err != nil {
       return
   }
   // Precompute time to cover distance d under constant acceleration a reaching no top speed
   tMax := math.Sqrt(float64(d*2) / float64(a))
   var ans float64
   for i := 0; i < n; i++ {
       var t, v int
       fmt.Fscan(reader, &t, &v)
       // time to accelerate to speed v and cover distance d
       dt := calc(float64(v), a, d, tMax)
       res := float64(t) + dt
       if res < ans {
           res = ans
       }
       fmt.Fprintf(writer, "%.5f\n", res)
       ans = res
   }
}

// calc returns time to cover distance d starting from speed 0,
// accelerating at rate a up to speed v, possibly cruising at v.
func calc(v float64, a, d int, tMax float64) float64 {
   t1 := v / float64(a)
   d1 := t1 * v / 2.0
   if d1 <= float64(d) {
       return (float64(d)-d1)/v + t1
   }
   return tMax
}
