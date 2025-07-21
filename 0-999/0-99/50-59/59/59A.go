package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   fmt.Fscan(reader, &s)
   upper, lower := 0, 0
   for i := 0; i < len(s); i++ {
       ch := s[i]
       if ch >= 'A' && ch <= 'Z' {
           upper++
       } else {
           lower++
       }
   }
   if upper > lower {
       fmt.Println(strings.ToUpper(s))
   } else {
       fmt.Println(strings.ToLower(s))
   }
}
