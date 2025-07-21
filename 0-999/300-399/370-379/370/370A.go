package main

import (
   "fmt"
)

func main() {
   var r1, c1, r2, c2 int
   if _, err := fmt.Scan(&r1, &c1, &r2, &c2); err != nil {
       return
   }
   // rook moves
   var rook int
   if r1 == r2 || c1 == c2 {
       rook = 1
   } else {
       rook = 2
   }
   // bishop moves
   var bishop int
   if (r1+c1)%2 != (r2+c2)%2 {
       bishop = 0
   } else if abs(r1-r2) == abs(c1-c2) {
       bishop = 1
   } else {
       bishop = 2
   }
   // king moves
   king := max(abs(r1-r2), abs(c1-c2))
   fmt.Printf("%d %d %d\n", rook, bishop, king)
}

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}
