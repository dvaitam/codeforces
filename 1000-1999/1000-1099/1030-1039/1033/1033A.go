package main

import (
   "fmt"
)

// reg computes a 2-bit region code for point (x,y) relative to the queen's position.
func reg(x, y, queenx, queeny int) int {
   var ans int
   ans <<= 1
   if x < queenx {
       ans += 1
   }
   ans <<= 1
   if y < queeny {
       ans += 1
   }
   return ans
}

func main() {
   var dims int
   var queenx, queeny int
   var bx, by, cx, cy int
   // Read board size (unused) and positions
   if _, err := fmt.Scan(&dims); err != nil {
       return
   }
   fmt.Scan(&queenx, &queeny)
   fmt.Scan(&bx, &by, &cx, &cy)

   if reg(bx, by, queenx, queeny) != reg(cx, cy, queenx, queeny) {
       fmt.Println("NO")
   } else {
       fmt.Println("YES")
   }
}
