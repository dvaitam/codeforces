package main

import (
   "fmt"
)

func main() {
   var a, b, c int64
   if _, err := fmt.Scan(&a, &b, &c); err != nil {
       return
   }
   // Let x = bonds between 1 and 2, y = between 2 and 3, z = between 3 and 1
   // Equations: x+z = a, x+y = b, y+z = c
   // Solve: x = (a + b - c) / 2, y = (b + c - a) / 2, z = (a + c - b) / 2
   sumABc := a + b - c
   sumBCa := b + c - a
   sumACb := a + c - b
   if sumABc < 0 || sumBCa < 0 || sumACb < 0 || sumABc%2 != 0 || sumBCa%2 != 0 || sumACb%2 != 0 {
       fmt.Println("Impossible")
       return
   }
   x := sumABc / 2
   y := sumBCa / 2
   z := sumACb / 2
   // final check: non-negative
   if x < 0 || y < 0 || z < 0 {
       fmt.Println("Impossible")
       return
   }
   fmt.Printf("%d %d %d", x, y, z)
}
