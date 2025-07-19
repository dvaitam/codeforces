package main

import (
   "fmt"
)

func main() {
   var n, k int64
   if _, err := fmt.Scan(&n, &k); err != nil {
       return
   }
   // compute ceil(k / n)
   res := (k + n - 1) / n
   fmt.Println(res)
}
