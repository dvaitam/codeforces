package main

import (
   "fmt"
)

func main() {
   var n, x1, y1, x2, y2 int
   if _, err := fmt.Scan(&n, &x1, &y1, &x2, &y2); err != nil {
       return
   }
   t1 := perimeterPosition(n, x1, y1)
   t2 := perimeterPosition(n, x2, y2)
   d := abs(t1 - t2)
   per := 4 * n
   if d > per-d {
       d = per - d
   }
   fmt.Println(d)
}

// perimeterPosition maps a point on the square boundary to a linear distance
// from the origin (0,0) going clockwise: bottom, right, top, left sides.
func perimeterPosition(n, x, y int) int {
   switch {
   case y == 0:
       return x
   case x == n:
       return n + y
   case y == n:
       return 2*n + (n - x)
   default: // x == 0
       return 3*n + (n - y)
   }
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}
