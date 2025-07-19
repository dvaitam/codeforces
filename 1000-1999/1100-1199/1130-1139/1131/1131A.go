package main

import (
   "fmt"
)

func min(x, y int) int {
   if x < y {
       return x
   }
   return y
}

func main() {
   var a, b, c, d int
   if _, err := fmt.Scan(&a, &b, &c, &d); err != nil {
       return
   }
   res := (a+2)*2 + (b+2)*2 - 4 + (c+2)*2 + (d+2)*2 - 4 - min(c+2, a+2)*2
   fmt.Println(res)
}
