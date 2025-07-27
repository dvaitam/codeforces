package main

import "fmt"

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   // remove optional sign
   if len(s) > 0 && (s[0] == '+' || s[0] == '-') {
       s = s[1:]
   }
   // output number of digits
   fmt.Println(len(s))
}
