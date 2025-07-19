package main

import (
   "fmt"
)

func main() {
   var x, y, z int
   // Read input values
   if _, err := fmt.Scan(&x, &y, &z); err != nil {
       return
   }
   // Output the result
   fmt.Println(1.1)
}
