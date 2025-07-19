package main

import (
   "fmt"
   "math"
)

// f computes the objective for given x, based on parameters A, B, L.
func f(x, A, B, L float64) float64 {
   y := math.Sqrt(L*L - x*x)
   return (A*x + B*y - x*y) / L
}

func main() {
   var A, B, L float64
   if _, err := fmt.Scanf("%f%f%f", &A, &B, &L); err != nil {
       return
   }
   // Handle edge cases
   if L <= B {
       fmt.Printf("%.8f\n", math.Min(A, L))
       return
   } else if L <= A {
       fmt.Printf("%.8f\n", math.Min(B, L))
       return
   }
   // Ternary-like search for maximum on [0, L]
   left, right := 0.0, L
   var ma, mb float64
   for i := 0; i < 500; i++ {
       ma = (left*3 + right) / 4.0
       mb = (left + right*3) / 4.0
       fa := f(ma, A, B, L)
       fb := f(mb, A, B, L)
       if fa < fb {
           right = mb
       } else {
           left = ma
       }
   }
   // Evaluate at midpoint
   ff := f((left+right)/2.0, A, B, L)
   ff = math.Min(ff, L)
   ff = math.Min(ff, A)
   if ff < 1e-8 {
       fmt.Println("My poor head =(")
   } else {
       fmt.Printf("%.8f\n", ff)
   }
}
