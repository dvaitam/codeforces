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
   delta int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var d float64
   if _, err := fmt.Fscan(in, &n, &d); err != nil {
       return
   }
   events := make([]event, 0, n*20*4)
   twoPi := 2 * math.Pi
   for i := 0; i < n; i++ {
       var xi, yi, ri float64
       fmt.Fscan(in, &xi, &yi, &ri)
       Di := math.Hypot(xi, yi)
       phi := math.Atan2(yi, xi)
       // possible k such that |k*d - Di| <= ri
       kmin := int(math.Ceil((Di - ri) / d))
       if kmin < 1 {
           kmin = 1
       }
       kmax := int(math.Floor((Di + ri) / d))
       for k := kmin; k <= kmax; k++ {
           // solve Di*cos(alpha-phi) in [kd - ri, kd + ri]
           kd := float64(k) * d
           A := (kd - ri) / Di
           B := (kd + ri) / Di
           if A > 1 || B < -1 {
               continue
           }
           // clamp
           if A < -1 {
               A = -1
           }
           if B > 1 {
               B = 1
           }
           // angles for cos(x) <= B and >= A
           a_l := math.Acos(A) // in [0, pi]
           a_u := math.Acos(B) // in [0, pi], a_l >= a_u
           // two intervals for x = alpha - phi
           // [a_u, a_l] and [2pi-a_l, 2pi-a_u]
           // map to alpha
           intervals := [][2]float64{
               {phi + a_u, phi + a_l},
               {phi + twoPi - a_l, phi + twoPi - a_u},
           }
           for _, iv := range intervals {
               s := iv[0]
               e := iv[1]
               // normalize to [0, 2pi)
               if s < 0 {
                   s += twoPi * (math.Floor(-s/twoPi) + 1)
               }
               if e < 0 {
                   e += twoPi * (math.Floor(-e/twoPi) + 1)
               }
               s = math.Mod(s, twoPi)
               e = math.Mod(e, twoPi)
               if e < s {
                   // wrap around
                   events = append(events, event{s, 1})
                   events = append(events, event{twoPi, -1})
                   events = append(events, event{0, 1})
                   events = append(events, event{e, -1})
               } else {
                   events = append(events, event{s, 1})
                   events = append(events, event{e, -1})
               }
           }
       }
   }
   // sweep line
   sort.Slice(events, func(i, j int) bool {
       if events[i].angle == events[j].angle {
           return events[i].delta > events[j].delta
       }
       return events[i].angle < events[j].angle
   })
   maxCnt := 0
   cur := 0
   for _, ev := range events {
       cur += ev.delta
       if cur > maxCnt {
           maxCnt = cur
       }
   }
   fmt.Println(maxCnt)
}
