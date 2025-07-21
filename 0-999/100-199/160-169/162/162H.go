package main

import (
   "fmt"
)

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   b := []byte(s)
   for i := range b {
       if i%2 == 0 {
           if b[i] >= 'a' && b[i] <= 'z' {
               b[i] = b[i] - 'a' + 'A'
           }
       } else {
           if b[i] >= 'A' && b[i] <= 'Z' {
               b[i] = b[i] - 'A' + 'a'
           }
       }
   }
   fmt.Println(string(b))
}
