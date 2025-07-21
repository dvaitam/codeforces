package main

import (
   "fmt"
)

func main() {
   var n, k, n1 int
   if _, err := fmt.Scan(&n, &k, &n1); err != nil {
       return
   }
   // If a single carpet covers whole room, or four carpets of size >= n/2 suffice
   if n1 >= n || k >= 4 {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
