package main

import (
   "fmt"
)

// floorDiv returns floor(u/m) for possibly negative u
func floorDiv(u, m int64) int64 {
   if u >= 0 {
       return u / m
   }
   // for negative u: (u+1)/m - 1 gives floor division
   return (u+1)/m - 1
}

func abs(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   var a, b, x1, y1, x2, y2 int64
   if _, err := fmt.Scan(&a, &b, &x1, &y1, &x2, &y2); err != nil {
       return
   }
   // Transform to u = x+y, v = x-y
   u1 := x1 + y1
   u2 := x2 + y2
   v1 := x1 - y1
   v2 := x2 - y2

   // Each bad line is u ≡ 0 mod 2a or v ≡ 0 mod 2b
   // Compute segment indices between consecutive bad lines
   iu1 := floorDiv(u1, 2*a)
   iu2 := floorDiv(u2, 2*a)
   iv1 := floorDiv(v1, 2*b)
   iv2 := floorDiv(v2, 2*b)

   // Number of required crossings for u-lines and v-lines
   pu := abs(iu1 - iu2)
   pv := abs(iv1 - iv2)
   // We can align crossings to minimize visits
   ans := max(pu, pv)
   fmt.Println(ans)
}
