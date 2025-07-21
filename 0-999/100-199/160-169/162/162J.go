package main

import (
   "fmt"
)

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   depth := 0
   for _, c := range s {
       if c == '(' {
           depth++
       } else {
           depth--
       }
       if depth < 0 {
           fmt.Println("NO")
           return
       }
   }
   if depth == 0 {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
