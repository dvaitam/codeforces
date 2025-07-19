package main

import (
   "fmt"
)

func main() {
   var k int
   if _, err := fmt.Scan(&k); err != nil {
       return
   }
   const add = 1 << 17
   // Output a 3x3 matrix where each row and column XOR equals k
   fmt.Println(3, 3)
   fmt.Println(k+add, k+add, k)
   fmt.Println(add, add, k+add)
   fmt.Println(0, 0, k)
}
