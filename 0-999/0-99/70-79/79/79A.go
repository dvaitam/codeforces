package main

import (
   "fmt"
)

func main() {
   var x, y int64
   if _, err := fmt.Scan(&x, &y); err != nil {
       return
   }
   // true: Ciel's turn; false: Hanako's turn
   turn := true
   for {
       if turn {
           // Ciel prefers max 100-yen coins => h=2,1,0
           // Fast-forward cycles where Ciel uses (2,2) and Hanako uses (0,22)
           if x >= 2 && y >= 24 {
               k1 := x / 2
               k2 := y / 24
               k := k1
               if k2 < k {
                   k = k2
               }
               if k > 0 {
                   x -= 2 * k
                   y -= 24 * k
                   continue
               }
           }
           // Single move
           if x >= 2 && y >= 2 {
               x -= 2; y -= 2
           } else if x >= 1 && y >= 12 {
               x -= 1; y -= 12
           } else if y >= 22 {
               y -= 22
           } else {
               fmt.Println("Hanako")
               return
           }
       } else {
           // Hanako prefers max 10-yen coins => h=0,1,2
           if y >= 22 {
               y -= 22
           } else if x >= 1 && y >= 12 {
               x -= 1; y -= 12
           } else if x >= 2 && y >= 2 {
               x -= 2; y -= 2
           } else {
               fmt.Println("Ciel")
               return
           }
       }
       turn = !turn
   }
}
