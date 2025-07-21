package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const (
   H  = 100.0
   Lx = 105.0
)

type Mirror struct {
   a, b float64
   v    int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var hl, hr float64
   var n int
   if _, err := fmt.Fscan(in, &hl, &hr, &n); err != nil {
       return
   }
   var floorMs, ceilMs []Mirror
   for i := 0; i < n; i++ {
       var vi int
       var ci byte
       var ai, bi float64
       fmt.Fscan(in, &vi, &ci, &ai, &bi)
       m := Mirror{a: ai, b: bi, v: vi}
       if ci == 'F' {
           floorMs = append(floorMs, m)
       } else {
           ceilMs = append(ceilMs, m)
       }
   }
   // sort by start
   sort.Slice(floorMs, func(i, j int) bool { return floorMs[i].a < floorMs[j].a })
   sort.Slice(ceilMs, func(i, j int) bool { return ceilMs[i].a < ceilMs[j].a })

   // helper to find mirror at boundary k
   findMirror := func(ms []Mirror, x float64) *Mirror {
       // binary search for largest a <= x
       i := sort.Search(len(ms), func(i int) bool { return ms[i].a > x })
       if i == 0 {
           return nil
       }
       m := &ms[i-1]
       if x+1e-9 >= m.a && x-1e-9 <= m.b {
           return m
       }
       return nil
   }

   maxScore := 0
   // direct no bounce
   // maxScore = 0 by default

   // two starting types: 0=F,1=T
   for s := 0; s < 2; s++ {
       // bounce count from 1 to n
       for mcnt := 1; mcnt <= n; mcnt++ {
           // generate boundary types b[k]: 0=F,1=T
           // compute image of hr by reflections reverse
           y := hr
           // reflect in order of bounces from last to first
           for j := mcnt; j >= 1; j-- {
               var btype int
               if j%2 == 1 {
                   btype = s
               } else {
                   btype = 1 - s
               }
               if btype == 0 {
                   // floor
                   y = -y
               } else {
                   // ceiling
                   y = 2*H - y
               }
           }
           dy := y - hl
           if dy == 0 {
               continue
           }
           score := 0
           ok := true
           // for each bounce
           for k := 1; k <= mcnt; k++ {
               // boundary type at k
               var btype int
               if k%2 == 1 {
                   btype = s
               } else {
                   btype = 1 - s
               }
               var planeY float64
               if dy > 0 {
                   planeY = float64(k) * H
               } else {
                   planeY = float64(k-1) * H
               }
               // x coordinate of intersection
               xk := Lx * (planeY - hl) / dy
               var m *Mirror
               if btype == 0 {
                   m = findMirror(floorMs, xk)
               } else {
                   m = findMirror(ceilMs, xk)
               }
               if m == nil {
                   ok = false
                   break
               }
               score += m.v
           }
           if ok && score > maxScore {
               maxScore = score
           }
       }
   }
   fmt.Println(maxScore)
}
