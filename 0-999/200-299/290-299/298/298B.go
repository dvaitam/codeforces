package main

import (
   "fmt"
   "os"
)

func main() {
   var t int
   var sx, sy, ex, ey int64
   // read input
   if _, err := fmt.Fscan(os.Stdin, &t, &sx, &sy, &ex, &ey); err != nil {
       return
   }
   var winds string
   if _, err := fmt.Fscan(os.Stdin, &winds); err != nil {
       return
   }
   // remaining distances to cover
   dx := ex - sx
   dy := ey - sy
   for i := 0; i < t; i++ {
       switch winds[i] {
       case 'E':
           if dx > 0 {
               dx--
           }
       case 'W':
           if dx < 0 {
               dx++
           }
       case 'N':
           if dy > 0 {
               dy--
           }
       case 'S':
           if dy < 0 {
               dy++
           }
       }
       if dx == 0 && dy == 0 {
           // 1-based time
           fmt.Println(i + 1)
           return
       }
   }
   // cannot reach
   fmt.Println(-1)
}
