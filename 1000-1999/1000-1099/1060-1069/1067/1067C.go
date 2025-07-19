package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   x := 0
   for i := 0; i < n; i++ {
       switch i % 3 {
       case 0, 1:
           x++
           fmt.Printf("%d %d\n", x, 0)
       case 2:
           fmt.Printf("%d %d\n", x, 3)
       }
   }
}
