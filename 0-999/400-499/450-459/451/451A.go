package main

import (
   "fmt"
)

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   // The game lasts for min(n, m) moves; Akshat wins if odd, otherwise Malvika wins.
   if n > m {
       n = m
   }
   if n%2 == 1 {
       fmt.Println("Akshat")
   } else {
       fmt.Println("Malvika")
   }
}
