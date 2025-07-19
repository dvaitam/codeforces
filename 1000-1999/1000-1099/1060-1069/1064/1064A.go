package main

import (
   "fmt"
)

func main() {
   var a, b, c int
   if _, err := fmt.Scan(&a, &b, &c); err != nil {
       return
   }
   // ensure a <= b <= c
   if a > b {
       a, b = b, a
   }
   if b > c {
       b, c = c, b
   }
   if a > b {
       a, b = b, a
   }
   // Check triangle inequality
   if a + b > c {
       fmt.Println(0)
   } else {
       // need to increase sticks to satisfy a + b > c
       fmt.Println(c - (a + b) + 1)
   }
}
