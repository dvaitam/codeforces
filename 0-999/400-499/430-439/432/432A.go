package main

import "fmt"

func main() {
   var n, k int
   if _, err := fmt.Scan(&n, &k); err != nil {
       return
   }
   valid := 0
   for i := 0; i < n; i++ {
       var y int
       fmt.Scan(&y)
       if y + k <= 5 {
           valid++
       }
   }
   // Each team consists of 3 people
   fmt.Println(valid / 3)
}
