package main

import (
   "fmt"
   "os"
)

func main() {
   var p int
   if _, err := fmt.Fscan(os.Stdin, &p); err != nil {
       return
   }
   for i := 0; i < 3; i++ {
       var a, b int
       if _, err := fmt.Fscan(os.Stdin, &a, &b); err != nil {
           return
       }
       if p == a {
           p = b
       } else if p == b {
           p = a
       }
   }
   fmt.Println(p)
}
