package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for t > 0 {
       t--
       var x, y int64
       var c [7]int64
       fmt.Fscan(reader, &x, &y)
       for i := 1; i <= 6; i++ {
           fmt.Fscan(reader, &c[i])
       }
       // minimize costs by relaxing adjacent directions
       for it := 0; it < 6; it++ {
           for i := 1; i <= 6; i++ {
               // neighbors in circular order
               left := i - 1
               if left == 0 {
                   left = 6
               }
               right := i + 1
               if right == 7 {
                   right = 1
               }
               if c[i] > c[left] + c[right] {
                   c[i] = c[left] + c[right]
               }
           }
       }
       var ans int64
       switch {
       case x >= 0 && y >= 0:
           // move with direction (1,1)
           k := x
           if y < k {
               k = y
           }
           ans += k * c[1]
           x -= k
           y -= k
           // remaining on axes
           if x > 0 {
               ans += x * c[6]
           } else if y > 0 {
               ans += y * c[2]
           }
       case x <= 0 && y <= 0:
           // both negative, move with (-1,-1)
           xx := -x
           yy := -y
           k := xx
           if yy < k {
               k = yy
           }
           ans += k * c[4]
           xx -= k
           yy -= k
           // remaining
           if xx > 0 {
               ans += xx * c[3]
           } else if yy > 0 {
               ans += yy * c[5]
           }
       case x >= 0 && y <= 0:
           // independent x and y
           ans += x * c[6]
           ans += -y * c[5]
       case x <= 0 && y >= 0:
           ans += y * c[2]
           ans += -x * c[3]
       }
       fmt.Println(ans)
   }
}
