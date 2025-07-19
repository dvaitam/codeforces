package main

import (
   "fmt"
)

func main() {
   var x, y int
   if _, err := fmt.Scanf("%d:%d", &x, &y); err != nil {
       return
   }
   for i := 0; i < 24; i++ {
       for j := 0; j < 60; j++ {
           if i > x || (i == x && j > y) {
               if i%10 == j/10 && i/10 == j%10 {
                   fmt.Printf("%02d:%02d\n", i, j)
                   return
               }
           }
       }
   }
   fmt.Println("00:00")
}
