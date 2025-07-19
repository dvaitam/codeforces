package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var ax, ay, bx, by, cx, cy int64
   if _, err := fmt.Fscan(reader, &ax, &ay, &bx, &by, &cx, &cy); err != nil {
       return
   }
   t := cx*cx + cy*cy
   // check if vector (x,y) can reach (bx,by) by integer multiples of (cx,cy)
   check := func(x, y int64) bool {
       if t == 0 {
           return x == bx && y == by
       }
       expr1 := cx*bx + cy*by - x*cx - y*cy
       expr2 := bx*cy - by*cx - x*cy + y*cx
       if expr1%t == 0 && expr2%t == 0 {
           return true
       }
       return false
   }
   // try original and rotated versions of (ax,ay)
   if check(ax, ay) || check(ay, -ax) || check(-ax, -ay) || check(-ay, ax) {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
