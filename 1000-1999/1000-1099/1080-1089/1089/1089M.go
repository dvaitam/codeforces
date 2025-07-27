package main

import (
   "fmt"
)

func main() {
   var a, b, c int
   if _, err := fmt.Scan(&a, &b, &c); err != nil {
       return
   }
   fmt.Println(a + b - c)
}
