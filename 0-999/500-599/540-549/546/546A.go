package main

import "fmt"

func main() {
   var k, n, w int
   fmt.Scan(&k, &n, &w)
   // total cost = k * (1 + 2 + ... + w)
   total := k * w * (w + 1) / 2
   // amount to borrow
   need := total - n
   if need < 0 {
       need = 0
   }
   fmt.Println(need)
}
