package main

import (
   "fmt"
)

func main() {
   var x, y, z int
   if _, err := fmt.Scan(&x, &y, &z); err != nil {
       return
   }
   // Compute minimum and maximum possible differences (upvotes - downvotes)
   dMin := x - y - z
   dMax := x + z - y
   switch {
   case dMin > 0:
       fmt.Println("+")
   case dMax < 0:
       fmt.Println("-")
   case dMin == 0 && dMax == 0:
       fmt.Println("0")
   default:
       fmt.Println("?")
   }
}
