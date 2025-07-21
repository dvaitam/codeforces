package main

import (
   "fmt"
)

func main() {
   var c1, c2, c3, c4 int
   if _, err := fmt.Scan(&c1, &c2, &c3, &c4); err != nil {
       return
   }
   var n, m int
   fmt.Scan(&n, &m)
   busCost := 0
   for i := 0; i < n; i++ {
       var a int
       fmt.Scan(&a)
       cost := a*c1
       if cost > c2 {
           cost = c2
       }
       busCost += cost
   }
   if busCost > c3 {
       busCost = c3
   }
   trolleyCost := 0
   for i := 0; i < m; i++ {
       var b int
       fmt.Scan(&b)
       cost := b*c1
       if cost > c2 {
           cost = c2
       }
       trolleyCost += cost
   }
   if trolleyCost > c3 {
       trolleyCost = c3
   }
   ans := busCost + trolleyCost
   if ans > c4 {
       ans = c4
   }
   fmt.Println(ans)
}
