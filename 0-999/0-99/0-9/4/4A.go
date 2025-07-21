package main

import (
   "fmt"
)

func main() {
   var w int
   if _, err := fmt.Scan(&w); err != nil {
       return
   }
   if w%2 == 0 && w > 2 {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
