package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   ans := -1
   maxv := -1
   for i := 0; i < n; i++ {
       var x int
       if _, err := fmt.Scan(&x); err != nil {
           return
       }
       if ans == -1 && x > maxv+1 {
           ans = i + 1
           break
       }
       if x > maxv {
           maxv = x
       }
   }
   fmt.Println(ans)
}
