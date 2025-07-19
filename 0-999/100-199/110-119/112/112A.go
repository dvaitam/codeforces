package main

import (
   "fmt"
   "strings"
)

func main() {
   var a, b string
   // read two strings
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   ua := strings.ToUpper(a)
   ub := strings.ToUpper(b)
   switch {
   case ua > ub:
       fmt.Print("1")
   case ua < ub:
       fmt.Print("-1")
   default:
       fmt.Print("0")
   }
}
