package main

import "fmt"

func main() {
   var a, b, c int
   // Read available fruits
   if _, err := fmt.Scan(&a); err != nil {
       return
   }
   fmt.Scan(&b)
   fmt.Scan(&c)
   // Determine maximum number of complete sets k
   k := a
   if b/2 < k {
       k = b / 2
   }
   if c/4 < k {
       k = c / 4
   }
   // Each set uses 1 lemon, 2 apples, 4 pears: total 7 fruits per set
   fmt.Println(7 * k)
}
