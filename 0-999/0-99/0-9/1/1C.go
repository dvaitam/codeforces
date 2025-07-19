package main

import (
   "fmt"
   "math"
)

const (
   eps = 1e-4
   PI  = math.Pi
)

func gcd(x, y float64) float64 {
   if y > eps {
       return gcd(y, math.Mod(x, y))
   }
   return x
}

func bcos(a, b, c float64) float64 {
   return math.Acos((a*a + b*b - c*c) / (2 * a * b))
}

func main() {
   var ax, ay, bx, by, cx, cy float64
   if _, err := fmt.Scanf("%f%f%f%f%f%f", &ax, &ay, &bx, &by, &cx, &cy); err != nil {
       return
   }
   a := math.Hypot(ax-bx, ay-by)
   b := math.Hypot(ax-cx, ay-cy)
   c := math.Hypot(bx-cx, by-cy)
   p := (a + b + c) / 2
   s := math.Sqrt(p * (p - a) * (p - b) * (p - c))
   R := (a * b * c) / (4 * s)
   A := bcos(b, c, a)
   B := bcos(a, c, b)
   C := bcos(a, b, c)
   n := PI / gcd(A, gcd(B, C))
   res := R*R*math.Sin(2*PI/n) * n / 2
   fmt.Printf("%.11f\n", res)
}
