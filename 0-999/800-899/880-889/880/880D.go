package main

import "fmt"

func main() {
   var a, b int
   fmt.Scan(&a, &b)
   if a < 0 {
       a = -a
   }
   if b < 0 {
       b = -b
   }
   for b != 0 {
       a, b = b, a%b
   }
   fmt.Println(a)
}
