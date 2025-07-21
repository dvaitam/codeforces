package main

import (
   "bufio"
   "fmt"
   "os"
)

func ceilDiv(a, b int64) int64 {
   if a >= 0 {
       return (a + b - 1) / b
   }
   // for negative a, floor division works as ceil
   return a / b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   ui := make([]int64, n)
   vi := make([]int64, n)
   var umin, umax, vmin, vmax int64
   for i := 0; i < n; i++ {
       var x, y int64
       fmt.Fscan(in, &x, &y)
       u := x + y
       v := x - y
       ui[i] = u
       vi[i] = v
       if i == 0 {
           umin, umax = u, u
           vmin, vmax = v, v
       } else {
           if u < umin {
               umin = u
           }
           if u > umax {
               umax = u
           }
           if v < vmin {
               vmin = v
           }
           if v > vmax {
               vmax = v
           }
       }
   }
   // direct-only optimum: half of max span in u or v
   tu := umax - umin
   tv := vmax - vmin
   // ceil(tu/2), ceil(tv/2)
   du := ceilDiv(tu, 2)
   dv := ceilDiv(tv, 2)
   directT := du
   if dv > directT {
       directT = dv
   }

   // if no stations, answer is directT
   if k == 0 {
       fmt.Println(directT)
       return
   }
   // read stations and compute transforms
   var m1, m2, M1, M2 int64
   for j := 0; j < k; j++ {
       var x, y int64
       fmt.Fscan(in, &x, &y)
       t1 := x + y
       t2 := x - y
       if j == 0 {
           m1, M1 = t1, t1
           m2, M2 = t2, t2
       } else {
           if t1 < m1 {
               m1 = t1
           }
           if t1 > M1 {
               M1 = t1
           }
           if t2 < m2 {
               m2 = t2
           }
           if t2 > M2 {
               M2 = t2
           }
       }
   }
   // compute max distance to nearest station
   var Amax int64
   for i := 0; i < n; i++ {
       u := ui[i]
       v := vi[i]
       // distance = max of four projections
       d1 := u - m1
       d2 := M1 - u
       d3 := v - m2
       d4 := M2 - v
       di := d1
       if d2 > di {
           di = d2
       }
       if d3 > di {
           di = d3
       }
       if d4 > di {
           di = d4
       }
       if di > Amax {
           Amax = di
       }
   }
   // answer is min(directT, Amax)
   ans := Amax
   if directT < ans {
       ans = directT
   }
   fmt.Println(ans)
}
