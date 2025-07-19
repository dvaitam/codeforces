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
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   type pt struct{ x, y, d float64 }
   pts := make([]pt, 0, n)
   c0 := 0
   for i := 0; i < n; i++ {
       var xi, yi int
       fmt.Fscan(in, &xi, &yi)
       if xi == 0 && yi == 0 {
           c0++
       } else {
           xf := float64(xi)
           yf := float64(yi)
           pts = append(pts, pt{xf, yf, math.Hypot(xf, yf)})
       }
   }
   if c0 >= k {
       fmt.Println("0.0")
       return
   }
   // search r in (0.5, 2e5+)
   l, r := 0.5, 222222.0
   // binary search on geometric mean
   for step := 0; step < 20; step++ {
       m := math.Sqrt(l * r)
       // build events
       type ev struct{ ang float64; v int }
       events := make([]ev, 0, len(pts)*4)
       for _, p := range pts {
           // check possible
           arg := p.d / (2 * m)
           if arg > 1 {
               continue
           }
           delta := math.Acos(arg)
           base := math.Atan2(p.y, p.x)
           a1 := base - delta
           a2 := base + delta
           events = append(events, ev{a1, 1}, ev{a2, -1}, ev{a1 + 2*math.Pi, 1}, ev{a2 + 2*math.Pi, -1})
       }
       sort.Slice(events, func(i, j int) bool {
           if events[i].ang != events[j].ang {
               return events[i].ang < events[j].ang
           }
           return events[i].v > events[j].v
       })
       run, hi := 0, 0
       for _, e := range events {
           run += e.v
           if run > hi {
               hi = run
           }
       }
       if hi >= k {
           r = m
       } else {
           l = m
       }
   }
   ans := math.Sqrt(l * r)
   // print with precision
   fmt.Printf("%.14f\n", ans)
}
