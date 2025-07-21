package main

import (
   "fmt"
   "strings"
)

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   var sb strings.Builder
   for i := 0; i < len(s); i++ {
       c := s[i]
       // skip vowels
       switch c {
       case 'A', 'O', 'Y', 'E', 'U', 'I', 'a', 'o', 'y', 'e', 'u', 'i':
           continue
       }
       // convert to lowercase if uppercase
       if c >= 'A' && c <= 'Z' {
           c = c - 'A' + 'a'
       }
       sb.WriteByte('.')
       sb.WriteByte(c)
   }
   fmt.Print(sb.String())
}
