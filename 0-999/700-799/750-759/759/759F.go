package main

import "fmt"

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   count := 0
   for _, c := range s {
       if c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u' {
           count++
       }
   }
   fmt.Println(count)
}
