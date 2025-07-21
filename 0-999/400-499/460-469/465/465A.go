package main

import (
   "fmt"
)

func main() {
   var n int
   var s string
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   ans := n
   for i := 0; i < n; i++ {
       if s[i] == '0' {
           ans = i + 1
           break
       }
   }
   fmt.Println(ans)
}
