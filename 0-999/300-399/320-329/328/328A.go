package main

import (
   "fmt"
)

func main() {
   // Read four integers
   var a [4]int64
   for i := 0; i < 4; i++ {
       if _, err := fmt.Scan(&a[i]); err != nil {
           return
       }
   }
   // Check arithmetic progression
   d1 := a[1] - a[0]
   d2 := a[2] - a[1]
   d3 := a[3] - a[2]
   if d1 == d2 && d2 == d3 {
       fmt.Println(a[3] + d1)
       return
   }
   // Check geometric progression (allow rational ratio)
   // Ensure non-zero first term
   if a[0] != 0 && a[1]*a[1] == a[0]*a[2] && a[2]*a[2] == a[1]*a[3] {
       num := a[1]
       den := a[0]
       // Next term = a[3] * num / den, must be integer
       if (a[3]*num)%den == 0 {
           fmt.Println(a[3] * num / den)
           return
       }
   }
   // Neither progression or next term not integer
   fmt.Println(42)
}
