package main

import (
   "fmt"
)

func main() {
   var rs, ks string
   if _, err := fmt.Scan(&rs); err != nil {
       return
   }
   if _, err := fmt.Scan(&ks); err != nil {
       return
   }
   // parse positions: file a-h -> 1-8, rank 1-8
   rX := int(rs[0]-'a') + 1
   rY := int(rs[1]-'0')
   kX := int(ks[0]-'a') + 1
   kY := int(ks[1]-'0')
   count := 0
   // iterate all squares for new knight
   for x := 1; x <= 8; x++ {
       for y := 1; y <= 8; y++ {
           // skip occupied
           if x == rX && y == rY {
               continue
           }
           if x == kX && y == kY {
               continue
           }
           // rook attacks along row/column
           if x == rX || y == rY {
               continue
           }
           // knight attack between new knight and rook
           if isKnight(rX, rY, x, y) {
               continue
           }
           // knight attack between new knight and existing knight
           if isKnight(kX, kY, x, y) {
               continue
           }
           count++
       }
   }
   fmt.Println(count)
}

// isKnight returns true if two positions are a knight's move apart
func isKnight(x1, y1, x2, y2 int) bool {
   dx := x1 - x2
   if dx < 0 {
       dx = -dx
   }
   dy := y1 - y2
   if dy < 0 {
       dy = -dy
   }
   return (dx == 1 && dy == 2) || (dx == 2 && dy == 1)
}
