package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   x1, x2, y1, y2 := 0, 1, 0, 1
   total := 4 + 3*n
   fmt.Println(total)
   fmt.Printf("%d %d\n", x1, y2)
   fmt.Printf("%d %d\n", x1, y1)
   fmt.Printf("%d %d\n", x2, y2)
   fmt.Printf("%d %d\n", x2, y1)
   for i := 0; i < n; i++ {
       x2++
       y2++
       x1++
       y1++
       fmt.Printf("%d %d\n", x1, y2)
       fmt.Printf("%d %d\n", x2, y2)
       fmt.Printf("%d %d\n", x2, y1)
   }
}
