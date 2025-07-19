package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   for i := 0; i < n; i++ {
       var x int
       fmt.Scan(&x)
       if x%2 == 0 {
           x--
       }
       if i > 0 {
           fmt.Print(" ")
       }
       fmt.Print(x)
   }
   // newline at end
   fmt.Println()
}
