package main

import (
   "fmt"
)

func main() {
   var n int
   // read the integer n
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   // output binary representation without leading zeros
   fmt.Printf("%b\n", n)
}
