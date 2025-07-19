package main

import (
   "fmt"
)

func main() {
   var x, y, z, t1, t2, t3 int64
   if _, err := fmt.Scan(&x, &y, &z, &t1, &t2, &t3); err != nil {
       return
   }
   stair := abs(x - y) * t1
   elev := (abs(x - z) + abs(y - x)) * t2 + 3*t3
   if elev <= stair {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}

func abs(a int64) int64 {
   if a < 0 {
       return -a
   }
   return a
}
