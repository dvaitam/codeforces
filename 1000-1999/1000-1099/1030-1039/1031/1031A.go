package main

import (
   "fmt"
)

func main() {
   var w, h, k int
   if _, err := fmt.Scan(&w, &h, &k); err != nil {
       return
   }
   res := 0
   for i := 0; i < k; i++ {
       sum := 2*h + 2*w - 4
       res += sum
       w -= 4
       h -= 4
   }
   fmt.Println(res)
}
