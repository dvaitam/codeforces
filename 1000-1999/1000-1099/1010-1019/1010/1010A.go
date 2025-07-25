package main

import (
   "fmt"
)

func main() {
   var n int
   var m float64
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   fmt.Scan(&m)
   a := make([]float64, n)
   b := make([]float64, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&a[i])
   }
   for i := 0; i < n; i++ {
       fmt.Scan(&b[i])
   }
   // If any coefficient equals 1, impossible (will consume infinite fuel)
   for i := 0; i < n; i++ {
       if a[i] == 1 || b[i] == 1 {
           fmt.Println(-1)
           return
       }
   }
   // simulate function for given fuel f
   can := func(f float64) bool {
       fuel := f
       for i := 0; i < n; i++ {
           // take-off from planet i
           w := m + fuel
           need := w / a[i]
           fuel -= need
           if fuel < 0 {
               return false
           }
           // landing to next planet
           w = m + fuel
           j := (i + 1) % n
           need = w / b[j]
           fuel -= need
           if fuel < 0 {
               return false
           }
       }
       return true
   }
   // binary search minimal f in [0, 1e9]
   lo, hi := 0.0, 1e9
   for it := 0; it < 200; it++ {
       mid := (lo + hi) / 2
       if can(mid) {
           hi = mid
       } else {
           lo = mid
       }
   }
   res := hi
   fmt.Printf("%.10f\n", res)
}
