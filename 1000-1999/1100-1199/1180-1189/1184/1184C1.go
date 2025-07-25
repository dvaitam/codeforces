package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   total := 4*n + 1
   type point struct{ x, y int }
   pts := make([]point, total)
   minX, maxX := int(1e9), int(-1e9)
   minY, maxY := int(1e9), int(-1e9)
   for i := 0; i < total; i++ {
       fmt.Fscan(in, &pts[i].x, &pts[i].y)
       if pts[i].x < minX {
           minX = pts[i].x
       }
       if pts[i].x > maxX {
           maxX = pts[i].x
       }
       if pts[i].y < minY {
           minY = pts[i].y
       }
       if pts[i].y > maxY {
           maxY = pts[i].y
       }
   }
   for _, p := range pts {
       if p.x > minX && p.x < maxX && p.y > minY && p.y < maxY {
           fmt.Println(p.x, p.y)
           return
       }
   }
}
