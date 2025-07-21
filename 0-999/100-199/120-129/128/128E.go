package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var k int64
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   xs := make([]float64, n)
   ys := make([]float64, n)
   rs := make([]float64, n)
   for i := 0; i < n; i++ {
       var xi, yi, ri float64
       fmt.Fscan(in, &xi, &yi, &ri)
       xs[i], ys[i], rs[i] = xi, yi, ri
   }
   // Compute max number of circles stabbed by a line
   M := 1
   pi := math.Pi
   for i := 0; i < n; i++ {
       // events: angle, delta
       type ev struct{ a float64; d int }
       events := make([]ev, 0, 2*(n-1))
       for j := 0; j < n; j++ {
           if j == i {
               continue
           }
           dx := xs[j] - xs[i]
           dy := ys[j] - ys[i]
           dist := math.Hypot(dx, dy)
           // disjoint circles, so dist > rs[j]
           // compute center angle
           phi := math.Atan2(dy, dx)
           // interval half-width
           ang := math.Asin(rs[j] / dist)
           width := ang * 2.0
           // raw left
           L := phi - ang
           // normalize L to [0, pi)
           L = math.Mod(L, pi)
           if L < 0 {
               L += pi
           }
           R := L + width
           if R < pi {
               events = append(events, ev{L, 1}, ev{R, -1})
           } else {
               // wrap
               events = append(events, ev{L, 1}, ev{pi, -1}, ev{0, 1}, ev{R - pi, -1})
           }
       }
       if len(events) == 0 {
           continue
       }
       // sort events by angle, start before end
       sort.Slice(events, func(a, b int) bool {
           if events[a].a == events[b].a {
               return events[a].d > events[b].d
           }
           return events[a].a < events[b].a
       })
       cnt, best := 0, 0
       for _, e := range events {
           cnt += e.d
           if cnt > best {
               best = cnt
           }
       }
       // include circle i itself
       if best+1 > M {
           M = best + 1
       }
   }
   // total banana pieces: n + k * M
   total := int64(n) + k*int64(M)
   fmt.Println(total)
}
