package main

import (
   "fmt"
)

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   count := 0
   for a := 0; a*a <= n; a++ {
       b := n - a*a
       if a + b*b == m {
           count++
       }
   }
   fmt.Println(count)
}
