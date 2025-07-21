package main

import (
   "fmt"
   "math"
)

func main() {
   var r, x1, y1, x2, y2 int64
   if _, err := fmt.Scan(&r, &x1, &y1, &x2, &y2); err != nil {
       return
   }
   dx := x1 - x2
   dy := y1 - y2
   dist := math.Hypot(float64(dx), float64(dy))
   maxMove := 2 * float64(r)
   steps := int64(math.Ceil(dist / maxMove))
   fmt.Println(steps)
}
