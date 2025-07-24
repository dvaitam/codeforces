package main

import (
   "fmt"
)

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   var n int
   var s string
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   // start of current segment (sx, sy)
   sx, sy := 0, 0
   // current position (x, y)
   x, y := 0, 0
   // number of segments (minimum possible m)
   count := 1
   // steps taken in current segment
   stepCount := 0
   for _, c := range s {
       // previous position before this move
       px, py := x, y
       switch c {
       case 'L':
           x--
       case 'R':
           x++
       case 'U':
           y++
       case 'D':
           y--
       }
       stepCount++
       // Manhattan distance from segment start to current
       dist := abs(x-sx) + abs(y-sy)
       // if path is no longer shortest, start new segment
       if dist < stepCount {
           count++
           // new segment starts from previous position
           sx, sy = px, py
           stepCount = 1
       }
   }
   fmt.Println(count)
}
