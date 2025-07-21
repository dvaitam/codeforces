package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var xa, ya, xb, yb int
   if _, err := fmt.Fscan(reader, &xa, &ya, &xb, &yb); err != nil {
       return
   }
   // normalize rectangle corners
   minX, maxX := xa, xb
   if minX > maxX {
       minX, maxX = maxX, minX
   }
   minY, maxY := ya, yb
   if minY > maxY {
       minY, maxY = maxY, minY
   }
   var n int
   fmt.Fscan(reader, &n)
   type heater struct{ x, y, r int }
   heaters := make([]heater, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &heaters[i].x, &heaters[i].y, &heaters[i].r)
   }
   covered := func(x, y int) bool {
       for _, h := range heaters {
           dx := x - h.x
           dy := y - h.y
           if dx*dx+dy*dy <= h.r*h.r {
               return true
           }
       }
       return false
   }
   blankets := 0
   // bottom and top edges
   for x := minX; x <= maxX; x++ {
       if !covered(x, minY) {
           blankets++
       }
       if minY != maxY && !covered(x, maxY) {
           blankets++
       }
   }
   // left and right edges (excluding corners)
   for y := minY + 1; y <= maxY-1; y++ {
       if !covered(minX, y) {
           blankets++
       }
       if minX != maxX && !covered(maxX, y) {
           blankets++
       }
   }
   fmt.Println(blankets)
}
