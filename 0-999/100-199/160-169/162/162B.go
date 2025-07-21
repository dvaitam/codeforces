package main

import (
   "fmt"
)

func main() {
   var n int64
   // Read input
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   // Output binary representation of n without leading zeros
   // %b formats as base-2
   fmt.Printf("%b", n)
}
