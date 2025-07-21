package main

import (
   "fmt"
)

func main() {
   var a, x, y int64
   if _, err := fmt.Scan(&a, &x, &y); err != nil {
       return
   }
   a2 := 2 * a
   x2 := 2 * x
   y2 := 2 * y
   // Check horizontal borders and below
   if y2 <= 0 || y2%a2 == 0 {
       fmt.Println(-1)
       return
   }
   // Determine row (1-based)
   row := y2/a2 + 1
   // Determine number of squares in row
   var w int
   if row == 1 || row%2 == 0 {
       w = 1
   } else {
       w = 2
   }
   // Determine position in row and check vertical borders
   var pos int
   switch w {
   case 1:
       // single square centered: x in (-a/2, a/2)
       // scaled boundaries at -a, a
       if x2 <= -a || x2 >= a {
           fmt.Println(-1)
           return
       }
       pos = 1
   case 2:
       // two squares: left [-a,0), right (0,a]
       // scaled boundaries at -2a,0,2a
       if x2 <= -a2 || x2 >= a2 || x2 == 0 {
           fmt.Println(-1)
           return
       }
       if x2 < 0 {
           pos = 1
       } else {
           pos = 2
       }
   }
   // Compute block number
   L := row - 1
   var sumPrev int64
   switch {
   case L <= 0:
       sumPrev = 0
   case L == 1:
       sumPrev = 1
   default:
       // width(1)=1; for k>=2: even->1, odd->2
       countEven := L / 2
       countOdd := (L+1)/2 - 1
       sumPrev = 1 + countEven + countOdd*2
   }
   fmt.Println(sumPrev + int64(pos))
}
