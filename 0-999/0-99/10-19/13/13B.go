package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

type Point struct { x, y int64 }

// check if C lies on segment AB (inclusive)
func onSegment(A, B, C Point) bool {
   // collinear: (B-A) x (C-A) == 0
   if (B.x-A.x)*(C.y-A.y) - (B.y-A.y)*(C.x-A.x) != 0 {
       return false
   }
   // dot((C-A),(C-B)) <= 0 ensures C between A and B
   if (C.x-A.x)*(C.x-B.x) + (C.y-A.y)*(C.y-B.y) > 0 {
       return false
   }
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       // read three segments
       seg := [3][2]Point{}
       for i := 0; i < 3; i++ {
           fmt.Fscan(reader,
               &seg[i][0].x, &seg[i][0].y,
               &seg[i][1].x, &seg[i][1].y)
       }
       ok := false
       // choose two segments as the legs
       for i := 0; i < 3 && !ok; i++ {
           for j := i + 1; j < 3 && !ok; j++ {
               // find common endpoint P
               var P, E1, E2 Point
               found := false
               a1, a2 := seg[i][0], seg[i][1]
               b1, b2 := seg[j][0], seg[j][1]
               // match endpoints
               if a1 == b1 {
                   P, E1, E2 = a1, a2, b2; found = true
               } else if a1 == b2 {
                   P, E1, E2 = a1, a2, b1; found = true
               } else if a2 == b1 {
                   P, E1, E2 = a2, a1, b2; found = true
               } else if a2 == b2 {
                   P, E1, E2 = a2, a1, b1; found = true
               }
               if !found {
                   continue
               }
               // angle check: 0 < angle <= 90
               v1x, v1y := float64(E1.x-P.x), float64(E1.y-P.y)
               v2x, v2y := float64(E2.x-P.x), float64(E2.y-P.y)
               dot := v1x*v2x + v1y*v2y
               l1 := math.Hypot(v1x, v1y)
               l2 := math.Hypot(v2x, v2y)
               if !(dot >= 0 && dot < l1*l2) {
                   continue
               }
               // third segment
               k := 3 - i - j
               C1, C2 := seg[k][0], seg[k][1]
               // try C1 on first leg and C2 on second, or vice versa
               tryPair := func(Q1, Q2 Point) bool {
                   if !onSegment(P, E1, Q1) || !onSegment(P, E2, Q2) {
                       return false
                   }
                   // check division ratios: t in [0.2,0.8]
                   // param dotParam1 = dot(Q1-P, E1-P), lensq1 = |E1-P|^2
                   dx1, dy1 := Q1.x-P.x, Q1.y-P.y
                   ex1, ey1 := E1.x-P.x, E1.y-P.y
                   dot1 := dx1*ex1 + dy1*ey1
                   lensq1 := ex1*ex1 + ey1*ey1
                   if dot1*5 < lensq1 || dot1*5 > 4*lensq1 {
                       return false
                   }
                   // second leg
                   dx2, dy2 := Q2.x-P.x, Q2.y-P.y
                   ex2, ey2 := E2.x-P.x, E2.y-P.y
                   dot2 := dx2*ex2 + dy2*ey2
                   lensq2 := ex2*ex2 + ey2*ey2
                   if dot2*5 < lensq2 || dot2*5 > 4*lensq2 {
                       return false
                   }
                   return true
               }
               if tryPair(C1, C2) || tryPair(C2, C1) {
                   ok = true
               }
           }
       }
       if ok {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
