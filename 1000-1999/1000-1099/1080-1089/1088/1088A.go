package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   if n == 1 {
       fmt.Print(-1)
   } else {
       fmt.Printf("%d %d", n, n)
   }
}
