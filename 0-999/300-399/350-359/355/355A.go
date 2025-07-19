package main

import (
   "fmt"
)

func main() {
   var k, d int
   if _, err := fmt.Scan(&k, &d); err != nil {
       return
   }
   // If more than one digit and sum zero: no solution
   if k >= 2 && d == 0 {
       fmt.Print("No solution")
       return
   }
   // Sum zero with single digit: zero
   if d == 0 {
       fmt.Print(0)
       return
   }
   // Single digit equals sum
   if k == 1 {
       fmt.Print(d)
       return
   }
   // For k >= 2 and d > 0, construct minimal number
   u := d - 1
   v := 1
   first := u
   second := v
   if v > u {
       first = v
       second = u
   }
   // Output first digit
   fmt.Print(first)
   // Middle zeros
   for i := 0; i < k-2; i++ {
       fmt.Print(0)
   }
   // Last digit
   fmt.Print(second)
}
