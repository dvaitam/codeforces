package main

import (
   "fmt"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   var a, b int
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   win1, draw, win2 := 0, 0, 0
   for x := 1; x <= 6; x++ {
       da := abs(a - x)
       db := abs(b - x)
       if da < db {
           win1++
       } else if da == db {
           draw++
       } else {
           win2++
       }
   }
   fmt.Println(win1, draw, win2)
}
