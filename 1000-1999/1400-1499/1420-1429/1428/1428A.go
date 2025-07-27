package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var x1, y1, x2, y2 int64
       fmt.Fscan(in, &x1, &y1, &x2, &y2)
       dx := x1 - x2
       if dx < 0 {
           dx = -dx
       }
       dy := y1 - y2
       if dy < 0 {
           dy = -dy
       }
       ans := dx + dy
       if dx > 0 && dy > 0 {
           ans += 2
       }
       fmt.Println(ans)
   }
}
