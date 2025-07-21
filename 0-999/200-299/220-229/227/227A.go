package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var xa, ya, xb, yb, xc, yc int64
   fmt.Fscan(reader, &xa, &ya)
   fmt.Fscan(reader, &xb, &yb)
   fmt.Fscan(reader, &xc, &yc)

   // vector from A to B
   vx := xb - xa
   vy := yb - ya
   // vector from B to C
   wx := xc - xb
   wy := yc - yb

   // cross product v x w = vx*wy - vy*wx
   cross := vx*wy - vy*wx

   switch {
   case cross > 0:
       fmt.Println("LEFT")
   case cross < 0:
       fmt.Println("RIGHT")
   default:
       fmt.Println("TOWARDS")
   }
}
