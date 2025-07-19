package main

import (
   "fmt"
)

func main() {
   var h, l float64
   // Read height h and diagonal l
   if _, err := fmt.Scan(&h, &l); err != nil {
       return
   }
   // Compute horizontal distance: (l^2 - h^2) / (2h)
   val := (l*l - h*h) / (2 * h)
   // Print with 10 decimal places
   fmt.Printf("%.10f\n", val)
}
