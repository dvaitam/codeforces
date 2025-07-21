package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   var x1, y1, x2, y2 int
   const INF = 1<<30
   minX, minY := INF, INF
   maxX, maxY := 0, 0
   var sumArea int64

   for i := 0; i < n; i++ {
       fmt.Scan(&x1, &y1, &x2, &y2)
       if x1 < minX {
           minX = x1
       }
       if y1 < minY {
           minY = y1
       }
       if x2 > maxX {
           maxX = x2
       }
       if y2 > maxY {
           maxY = y2
       }
       w := x2 - x1
       h := y2 - y1
       sumArea += int64(w) * int64(h)
   }

   dx := maxX - minX
   dy := maxY - minY
   if dx != dy {
       fmt.Println("NO")
       return
   }
   if sumArea != int64(dx)*int64(dy) {
       fmt.Println("NO")
       return
   }
   fmt.Println("YES")
}
