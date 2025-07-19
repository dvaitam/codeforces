package main

import (
   "fmt"
   "math"
   "sort"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   radii := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&radii[i])
   }
   sort.Ints(radii)
   var s float64
   for i, r := range radii {
       area := math.Pi * float64(r) * float64(r)
       if n%2 == 1 {
           if i%2 == 0 {
               s += area
           } else {
               s -= area
           }
       } else {
           if i%2 == 0 {
               s -= area
           } else {
               s += area
           }
       }
   }
   fmt.Printf("%.15f\n", s)
}
