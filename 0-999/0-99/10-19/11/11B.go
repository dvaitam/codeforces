package main

import "fmt"

func main() {
   var x int
   // read input
   if _, err := fmt.Scan(&x); err != nil {
       return
   }
   if x < 0 {
       x = -x
   }
   k, p := 0, 0
   // find smallest k such that sum 1..k >= x and parity matches x
   for p < x || (x%2 != p%2) {
       k++
       p += k
   }
   fmt.Println(k)
}
