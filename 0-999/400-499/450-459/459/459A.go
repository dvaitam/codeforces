package main

import (
   "fmt"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   var a, b, c, d int
   if _, err := fmt.Scan(&a, &b, &c, &d); err != nil {
       return
   }
   // Check for vertical alignment
   if a == c {
       w := abs(b - d)
       x3, y3 := a - w, b
       x4, y4 := c - w, d
       fmt.Printf("%d %d %d %d\n", x3, y3, x4, y4)
       return
   }
   // Check for horizontal alignment
   if b == d {
       w := abs(a - c)
       x3, y3 := a, b - w
       x4, y4 := c, d - w
       fmt.Printf("%d %d %d %d\n", x3, y3, x4, y4)
       return
   }
   // Check for diagonal of square
   if abs(a-c) == abs(b-d) {
       // The other two vertices
       fmt.Printf("%d %d %d %d\n", a, d, c, b)
       return
   }
   // No solution
   fmt.Println(-1)
}
