package main

import (
   "fmt"
)

func main() {
   var a, b, c int
   if _, err := fmt.Scan(&a, &b, &c); err != nil {
       return
   }
   // total number of hexagonal tiles in hexagon with sides a,b,c,a,b,c
   result := a*b + b*c + a*c - a - b - c + 1
   fmt.Println(result)
}
