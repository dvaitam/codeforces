package main

import (
   "fmt"
)

func main() {
   var n, s int
   if _, err := fmt.Scan(&n, &s); err != nil {
       return
   }
   ans := -1
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Scan(&x, &y)
       // total price in cents
       cost := x*100 + y
       if s*100 >= cost {
           // sweets is the cents part of the change
           sweets := (100 - y) % 100
           if sweets > ans {
               ans = sweets
           }
       }
   }
   fmt.Println(ans)
}
