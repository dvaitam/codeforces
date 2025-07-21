package main

import (
   "fmt"
   "unicode"
)

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   r := []rune(s)
   if len(r) > 0 {
       r[0] = unicode.ToUpper(r[0])
   }
   fmt.Println(string(r))
}
