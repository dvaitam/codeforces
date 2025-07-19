package main

import (
   "fmt"
)

func main() {
   var a int
   if _, err := fmt.Scan(&a); err != nil {
       return
   }
   // Set of winning positions
   winning := map[int]bool{
       2:  true,
       3:  true,
       4:  true,
       5:  true,
       12: true,
       30: true,
       31: true,
       35: true,
       43: true,
       46: true,
       52: true,
       64: true,
       86: true,
   }
   if winning[a] {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
