package main

import (
   "fmt"
)

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   count := 1
   for i := 1; i < len(s); i++ {
       if s[i] == s[i-1] {
           count++
           if count >= 7 {
               fmt.Println("YES")
               return
           }
       } else {
           count = 1
       }
   }
   fmt.Println("NO")
}
