package main

import (
   "fmt"
   "math"
)

func main() {
   var a, b, c float64
   // Read coefficients
   fmt.Scan(&a, &b, &c)
   // No solutions
   if a == 0 && b == 0 && c != 0 {
       fmt.Print("0")
       return
   }
   // Infinite solutions
   if a == 0 && b == 0 && c == 0 {
       fmt.Print("-1")
       return
   }
   // Linear equation
   if a == 0 {
       r := -c / b
       fmt.Println("1")
       fmt.Printf("%f", r)
       return
   }
   // Quadratic equation
   disc := b*b - 4*a*c
   if disc < 0 {
       fmt.Print("0")
       return
   }
   d := math.Sqrt(disc)
   // One real root
   if d == 0 {
       r := -b / (2 * a)
       fmt.Println("1")
       fmt.Printf("%f", r)
       return
   }
   // Two real roots
   r1 := (-b - d) / (2 * a)
   r2 := (-b + d) / (2 * a)
   fmt.Println("2")
   if r1 > r2 {
       fmt.Printf("%f\n%f\n", r2, r1)
   } else {
       fmt.Printf("%f\n%f\n", r1, r2)
   }
}
