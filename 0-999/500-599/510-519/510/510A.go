package main

import (
   "fmt"
   "strings"
)

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   for i := 1; i <= n; i++ {
       switch {
       case i%2 == 1:
           fmt.Println(strings.Repeat("#", m))
       case i%4 == 2:
           fmt.Println(strings.Repeat(".", m-1) + "#")
       case i%4 == 0:
           fmt.Println("#" + strings.Repeat(".", m-1))
       }
   }
}
