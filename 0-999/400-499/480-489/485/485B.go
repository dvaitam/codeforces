package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var x, y int64
   // Initialize min and max with first point
   // Read first point separately
   if _, err := fmt.Fscan(reader, &x, &y); err != nil {
       return
   }
   minX, maxX := x, x
   minY, maxY := y, y
   for i := 1; i < n; i++ {
       if _, err := fmt.Fscan(reader, &x, &y); err != nil {
           break
       }
       if x < minX {
           minX = x
       }
       if x > maxX {
           maxX = x
       }
       if y < minY {
           minY = y
       }
       if y > maxY {
           maxY = y
       }
   }
   dx := maxX - minX
   dy := maxY - minY
   side := dx
   if dy > side {
       side = dy
   }
   // Compute area: side squared
   area := side * side
   fmt.Println(area)
}
