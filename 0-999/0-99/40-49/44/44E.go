package main

import "fmt"

func main() {
   var k, a, b int
   var str string
   if _, err := fmt.Scan(&k, &a, &b); err != nil {
       return
   }
   if _, err := fmt.Scan(&str); err != nil {
       return
   }
   n := len(str)
   if n < a*k || n > b*k {
       fmt.Println("No solution")
       return
   }
   ix := 0
   for k > 0 {
       segLen := (n - ix) / k
       fmt.Println(str[ix : ix+segLen])
       ix += segLen
       k--
   }
}
