package main

import (
   "fmt"
)

func main() {
   var n, a, x, b, y int
   if _, err := fmt.Scan(&n, &a, &x, &b, &y); err != nil {
       return
   }
   posD, posV := a, b
   for {
       if posD == posV {
           fmt.Println("YES")
           return
       }
       if posD == x || posV == y {
           fmt.Println("NO")
           return
       }
       posD++
       if posD > n {
           posD = 1
       }
       posV--
       if posV < 1 {
           posV = n
       }
   }
}
