package main

import "fmt"

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   balance := 0
   for _, ch := range s {
       if ch == '(' {
           balance++
       } else if ch == ')' {
           balance--
       }
       if balance < 0 {
           fmt.Println("NO")
           return
       }
   }
   if balance == 0 {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
