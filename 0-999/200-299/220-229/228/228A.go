package main

import "fmt"

func main() {
   var a, b, c, d int
   // Read four horseshoe colors
   if _, err := fmt.Scan(&a, &b, &c, &d); err != nil {
       return
   }
   // Count distinct colors
   set := make(map[int]bool)
   set[a] = true
   set[b] = true
   set[c] = true
   set[d] = true
   // Minimum to buy is 4 minus distinct count
   fmt.Println(4 - len(set))
}
