package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   div := 1
   q := n
   for i := 9; i >= 1; i-- {
       if n%i == 0 {
           div = i
           q = n / i
           break
       }
   }
   fmt.Println(q)
   for j := 0; j < q; j++ {
       if j > 0 {
           fmt.Print(" ")
       }
       fmt.Print(div)
   }
   fmt.Println()
}
