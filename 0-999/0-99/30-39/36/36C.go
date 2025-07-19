package main

import (
   "fmt"
   "os"
)

// heightIntersect computes intersection height between two truncated cones
func heightIntersect(h0, r0, R0, h1, r1, R1 float64) float64 {
   res1 := (r1 - r0) / (R0 - r0) * h0
   res2 := (R0 - r1) / (R1 - r1) * h1
   res3 := (R1 - r0) / (R0 - r0) * h0
   if res2 > h1 {
       res2 = -1000
   } else {
       res2 = h0 - res2
   }
   if res3 > h0 {
       res3 = h0
   }
   if res3 < 0 {
       res3 = 0
   }
   res3 -= h1
   if r1 >= R0 {
       return h0
   }
   if res1 < 0 {
       res1 = 0
   }
   if res2 > h0 {
       res2 = h0
   }
   if res1 > res2 && res1 > res3 {
       return res1
   }
   if res3 > res2 {
       return res3
   }
   return res2
}

func main() {
   var n int
   fmt.Fscan(os.Stdin, &n)
   h := make([]float64, n)
   r := make([]float64, n)
   R := make([]float64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(os.Stdin, &h[i], &r[i], &R[i])
   }
   // DP state arrays
   sh := make([]float64, n)
   ch := make([]float64, n)
   cr := make([]float64, n)
   cR := make([]float64, n)
   curr := 1
   sh[0] = 0
   ch[0] = h[0]
   cr[0] = r[0]
   cR[0] = R[0]
   for i := 1; i < n; i++ {
       mh := 0.0
       for j := 0; j < curr; j++ {
           hh := heightIntersect(ch[j], cr[j], cR[j], h[i], r[i], R[i]) + sh[j]
           if hh > mh {
               mh = hh
           }
       }
       prog := 0
       for j := 0; j < curr; j++ {
           if sh[j]+ch[j] > mh {
               sh[prog] = sh[j]
               ch[prog] = ch[j]
               cr[prog] = cr[j]
               cR[prog] = cR[j]
               prog++
           }
       }
       sh[prog] = mh
       ch[prog] = h[i]
       cr[prog] = r[i]
       cR[prog] = R[i]
       curr = prog + 1
   }
   res := 0.0
   for i := 0; i < curr; i++ {
       if sh[i]+ch[i] > res {
           res = sh[i] + ch[i]
       }
   }
   fmt.Printf("%.12f\n", res)
}
