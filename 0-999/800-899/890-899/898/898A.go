package main

import (
   "fmt"
)

func main() {
   var n int
   _, err := fmt.Scan(&n)
   if err != nil {
       return
   }
   c := n % 10
   n -= c
   if c >= 5 {
       n += 10
   }
   fmt.Println(n)
}
