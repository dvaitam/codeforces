package main

import (
   "fmt"
   "math"
)

func main() {
   var hh, mm int
   var H, D, C, N int
   if _, err := fmt.Scanf("%d %d", &hh, &mm); err != nil {
       return
   }
   if _, err := fmt.Scanf("%d %d %d %d", &H, &D, &C, &N); err != nil {
       return
   }
   // Calculate burgers needed now
   burgers := int(math.Ceil(float64(H) / float64(N)))
   // Cost if buy now at full price
   costNow := float64(burgers) * float64(C)
   // If after or at discount time (20:00), apply discount immediately
   if hh >= 20 {
       discounted := float64(burgers) * (float64(C) * 4.0 / 5.0)
       fmt.Printf("%.5f", discounted)
       return
   }
   // Compute additional fat until 20:00
   minutesUntil := (20-hh)*60 - mm
   H += minutesUntil * D
   // Recompute burgers and cost with discount
   burgers2 := int(math.Ceil(float64(H) / float64(N)))
   costDiscount := float64(burgers2) * (float64(C) * 4.0 / 5.0)
   // Output minimum cost
   if costNow < costDiscount {
       fmt.Printf("%.5f", costNow)
   } else {
       fmt.Printf("%.5f", costDiscount)
   }
}
