package main

import (
   "fmt"
)

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   n := len(s)
   x := -1
   for i := 1; i < n; i++ {
       if s[i-1] == 'r' && s[i] == 'u' {
           x = i - 1
       }
   }
   start := 0
   if n > 0 && s[0] == 'h' {
       fmt.Print("http://")
       start = 4
   } else {
       fmt.Print("ftp://")
       start = 3
   }
   if x > start {
       fmt.Print(s[start:x])
   }
   fmt.Print(".ru")
   if x+2 < n {
       fmt.Print("/")
       fmt.Print(s[x+2:])
   }
}
