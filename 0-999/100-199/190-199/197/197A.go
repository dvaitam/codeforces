package main

import (
   "fmt"
)

func main() {
   var a, b, r int
   if _, err := fmt.Scan(&a, &b, &r); err != nil {
       return
   }
   // Number of plates along each dimension
   na := a / (2 * r)
   nb := b / (2 * r)
   total := na * nb
   if total%2 == 1 {
       fmt.Println("First")
   } else {
       fmt.Println("Second")
   }
}
