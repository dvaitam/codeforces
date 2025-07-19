package main

import (
   "fmt"
   "math"
)

func main() {
   var k, d, t float64
   if _, err := fmt.Scan(&k, &d, &t); err != nil {
       return
   }
   // Determine cycle length: smallest multiple of d not less than k
   var cycle float64
   if k > d {
       cycle = math.Ceil(k/d) * d
   } else {
       cycle = d
   }
   // Scale consumption: work phase at rate 2, idle at rate 1
   totalNeed := 2 * t
   perCycle := k + cycle
   // Full cycles
   full := math.Floor(totalNeed / perCycle)
   rem := totalNeed - perCycle*full
   var extra float64
   if rem <= 2*k {
       extra = rem / 2
   } else {
       extra = k + (rem-2*k)
   }
   // Total time = full cycles * cycle length + extra time
   res := full*cycle + extra
   // Print with sufficient precision
   fmt.Printf("%.10f\n", res)
}
