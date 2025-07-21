package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var targetX, targetY int
   if _, err := fmt.Fscan(reader, &targetX, &targetY); err != nil {
       return
   }
   // If target is starting point, no turns needed
   if targetX == 0 && targetY == 0 {
       fmt.Println(0)
       return
   }
   // Directions: right, up, left, down
   dirs := [4][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
   x, y := 0, 0
   d := 0
   turns := 0
   first := true
   // Spiral with segment lengths: 1,1,2,2,3,3,...
   for k := 1; ; k++ {
       for rep := 0; rep < 2; rep++ {
           if first {
               first = false
           } else {
               d = (d + 1) % 4
               turns++
           }
           dx, dy := dirs[d][0], dirs[d][1]
           for step := 0; step < k; step++ {
               x += dx
               y += dy
               if x == targetX && y == targetY {
                   fmt.Println(turns)
                   return
               }
           }
       }
   }
}
