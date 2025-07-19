package main

import (
   "fmt"
)

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   // simulate removals of consecutive equal letters
   stack := make([]rune, 0, len(s))
   cnt := 0
   for _, ch := range s {
       n := len(stack)
       if n > 0 && stack[n-1] == ch {
           // remove matching pair
           stack = stack[:n-1]
           cnt++
       } else {
           stack = append(stack, ch)
       }
   }
   // first player wins if number of removals is odd
   if cnt%2 == 1 {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
