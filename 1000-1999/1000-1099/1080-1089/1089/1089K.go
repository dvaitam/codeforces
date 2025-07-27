package main

import (
   "fmt"
)

func main() {
   var a, b, c int
   if _, err := fmt.Scan(&a, &b, &c); err != nil {
       return
   }
   // sort three values
   if a > b {
       a, b = b, a
   }
   if b > c {
       b, c = c, b
   }
   if a > b {
       a, b = b, a
   }
   fmt.Println(a, b, c)
}
