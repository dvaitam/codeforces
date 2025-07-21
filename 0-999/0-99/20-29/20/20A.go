package main

import (
   "fmt"
   "strings"
)

func main() {
   var s string
   // Read the path
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   var b strings.Builder
   prevSlash := false
   // Collapse multiple slashes
   for _, ch := range s {
       if ch == '/' {
           if !prevSlash {
               b.WriteRune('/')
               prevSlash = true
           }
       } else {
           b.WriteRune(ch)
           prevSlash = false
       }
   }
   res := b.String()
   // Remove trailing slash unless it's the root
   if len(res) > 1 && res[len(res)-1] == '/' {
       res = res[:len(res)-1]
   }
   fmt.Println(res)
}
