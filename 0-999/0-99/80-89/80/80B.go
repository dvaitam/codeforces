package main

import (
   "fmt"
)

func main() {
   var h, m int
   // Read time in H:M format
   if _, err := fmt.Scanf("%d:%d", &h, &m); err != nil {
       return
   }
   // Calculate angles
   hourAngle := float64(h%12)*30.0 + float64(m)/2.0
   minuteAngle := float64(m) * 6.0
   // Print with 9 decimal places
   fmt.Printf("%.9f %.9f\n", hourAngle, minuteAngle)
}
