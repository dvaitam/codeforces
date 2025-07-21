package main

import (
   "fmt"
)

func main() {
   var n int64
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   // pentagonal number formula: p(n) = (3n^2 - n) / 2
   result := (3*n*n - n) / 2
   fmt.Println(result)
}
