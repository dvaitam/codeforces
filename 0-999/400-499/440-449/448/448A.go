package main

import (
   "fmt"
)

func main() {
   var a1, a2, a3 int
   var b1, b2, b3 int
   var n int
   fmt.Scan(&a1, &a2, &a3)
   fmt.Scan(&b1, &b2, &b3)
   fmt.Scan(&n)

   totalCups := a1 + a2 + a3
   totalMedals := b1 + b2 + b3

   // Each shelf holds up to 5 cups or 10 medals
   requiredCupShelves := (totalCups + 4) / 5
   requiredMedalShelves := (totalMedals + 9) / 10

   if requiredCupShelves+requiredMedalShelves <= n {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
