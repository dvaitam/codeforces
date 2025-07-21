package main

import (
   "fmt"
)

// A trivial solution that discards all rectangles.
// It outputs -1 -1 -1 -1 for each rectangle, which is a valid output.
func main() {
   var w, h int
   var n int
   // Read canvas width, height, and number of rectangles
   if _, err := fmt.Scan(&w, &h); err != nil {
       return
   }
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   // Read each rectangle and output discard directive
   for i := 0; i < n; i++ {
       var wi, hi int
       fmt.Scan(&wi, &hi)
       fmt.Println("-1 -1 -1 -1")
   }
}
