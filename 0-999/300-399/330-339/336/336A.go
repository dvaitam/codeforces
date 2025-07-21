package main

import (
   "fmt"
)

func abs64(a int64) int64 {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   var x, y int64
   if _, err := fmt.Scan(&x, &y); err != nil {
       return
   }
   // Length of legs
   L := abs64(x) + abs64(y)
   // Points on axes
   // P1 on x-axis, P2 on y-axis
   x1 := (func() int64 { if x > 0 { return L } return -L })()
   y1 := int64(0)
   x2 := int64(0)
   y2 := (func() int64 { if y > 0 { return L } return -L })()
   // ensure x1 < x2
   if x1 >= x2 {
       // swap P1 and P2
       x1, x2 = x2, x1
       y1, y2 = y2, y1
   }
   fmt.Printf("%d %d %d %d\n", x1, y1, x2, y2)
}
