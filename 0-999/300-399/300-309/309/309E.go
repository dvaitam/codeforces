package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   sum := 0
   for i := 0; i < n; i++ {
       var x int
       if _, err := fmt.Scan(&x); err != nil {
           return
       }
       sum += x
   }
   fmt.Println(sum)
}
