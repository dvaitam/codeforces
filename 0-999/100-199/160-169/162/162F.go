package main

import (
   "fmt"
)

func main() {
   var n int64
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   var count int64
   for p := int64(5); p <= n; p *= 5 {
       count += n / p
   }
   fmt.Println(count)
}
