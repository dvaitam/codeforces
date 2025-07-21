package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, a, b int64
   if _, err := fmt.Fscan(reader, &n, &a, &b); err != nil {
       return
   }
   // required minimum area
   m := 6 * n
   // if current room is enough
   if a*b >= m {
       fmt.Printf("%d %d %d", a*b, a, b)
       return
   }
   swapped := false
   // ensure a <= b
   if a > b {
       a, b = b, a
       swapped = true
   }
   bestArea := int64(1)<<62 - 1
   var bestA, bestB int64
   // limit for iteration
   lim := int64(math.Sqrt(float64(m))) + 1
   for A := a; A <= lim; A++ {
       // minimal B to satisfy area
       B := (m + A - 1) / A
       if B < b {
           B = b
       }
       area := A * B
       if area < bestArea {
           bestArea = area
           bestA, bestB = A, B
       }
   }
   // consider using B = b and adjust A
   A := (m + b - 1) / b
   if A < a {
       A = a
   }
   area := A * b
   if area < bestArea {
       bestArea = area
       bestA, bestB = A, b
   }
   // swap back if needed
   if swapped {
       bestA, bestB = bestB, bestA
   }
   fmt.Printf("%d %d %d", bestArea, bestA, bestB)
}
