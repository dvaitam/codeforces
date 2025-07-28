package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

type event struct {
   angle float64
   typ   int // 0=T,1=R
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   pts := make([][2]float64, n)
   for i := 0; i < n; i++ {
       var x, y int64
       fmt.Fscan(in, &x, &y)
       pts[i][0] = float64(x)
       pts[i][1] = float64(y)
   }
   // initial support pointers at theta=0
   r, l, t, b := 0, 0, 0, 0
   for i := 1; i < n; i++ {
       if pts[i][0] > pts[r][0] {
           r = i
       }
       if pts[i][0] < pts[l][0] {
           l = i
       }
       if pts[i][1] > pts[t][1] {
           t = i
       }
       if pts[i][1] < pts[b][1] {
           b = i
       }
   }
   // build events
   events := make([]event, 0, 2*n)
   pi := math.Pi
   halfPi := pi / 2
   for i := 0; i < n; i++ {
       j := (i + 1) % n
       dx := pts[j][0] - pts[i][0]
       dy := pts[j][1] - pts[i][1]
       phi := math.Atan2(dy, dx)
       // normalize to [0,pi)
       if phi < 0 {
           phi += pi
       }
       if phi >= pi {
           phi -= pi
       }
       // T event at phi
       events = append(events, event{angle: phi, typ: 0})
       // R event at phi+pi/2 mod pi
       a := phi + halfPi
       if a >= pi {
           a -= pi
       }
       events = append(events, event{angle: a, typ: 1})
   }
   // sort events by angle
   sort.Slice(events, func(i, j int) bool {
       return events[i].angle < events[j].angle
   })
   // process intervals
   total := 0.0
   prev := 0.0
   // apply events at angle==0
   idx := 0
   for idx < len(events) && events[idx].angle == 0 {
       if events[idx].typ == 0 {
           t = (t + 1) % n
           b = (b - 1 + n) % n
       } else {
           r = (r + 1) % n
           l = (l - 1 + n) % n
       }
       idx++
   }
   // iterate events
   for ; idx < len(events); idx++ {
       ang := events[idx].angle
       if ang > prev {
           total += integrateInterval(prev, ang, pts, r, l, t, b)
           prev = ang
       }
       // update pointer
       if events[idx].typ == 0 {
           t = (t + 1) % n
           b = (b - 1 + n) % n
       } else {
           r = (r + 1) % n
           l = (l - 1 + n) % n
       }
   }
   // last interval
   if prev < pi {
       total += integrateInterval(prev, pi, pts, r, l, t, b)
   }
   // expected value over [0,pi)
   fmt.Printf("%.9f\n", total/pi)
}

func integrateInterval(a, b float64, pts [][2]float64, r, l, t, btm int) float64 {
   // width coeffs C, D
   Cr := pts[r][0] - pts[l][0]
   Dr := pts[r][1] - pts[l][1]
   // height coeffs A sin + B cos, where A = -(x_t-x_b), B=y_t-y_b
   At := -(pts[t][0] - pts[btm][0])
   Bt := pts[t][1] - pts[btm][1]
   L := b - a
   sin2b := math.Sin(2 * b)
   sin2a := math.Sin(2 * a)
   cos2b := math.Cos(2 * b)
   cos2a := math.Cos(2 * a)
   Icos2 := L/2 + (sin2b - sin2a)/4
   Isin2 := L/2 - (sin2b - sin2a)/4
   IsinCos := (-cos2b + cos2a) / 4
   return Cr*Bt*Icos2 + Dr*At*Isin2 + (Dr*Bt+Cr*At)*IsinCos
}
