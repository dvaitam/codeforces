package main

import (
   "fmt"
)

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   buses := 0
   current := 0
   for i := 0; i < n; i++ {
       var a int
       fmt.Scan(&a)
       if current+a <= m {
           current += a
       } else {
           buses++
           current = a
       }
   }
   if current > 0 {
       buses++
   }
   fmt.Println(buses)
}
