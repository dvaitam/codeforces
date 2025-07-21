package main

import "fmt"

func main() {
   var p string
   if _, err := fmt.Scan(&p); err != nil {
       return
   }
   for _, c := range p {
       if c == 'H' || c == 'Q' || c == '9' {
           fmt.Println("YES")
           return
       }
   }
   fmt.Println("NO")
}
