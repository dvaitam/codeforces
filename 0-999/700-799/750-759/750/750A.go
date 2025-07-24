package main

import "fmt"

func main() {
   var n, k int
   fmt.Scan(&n, &k)
   rem := 240 - k
   solved := 0
   timeSpent := 0
   for i := 1; i <= n; i++ {
       if timeSpent+5*i > rem {
           break
       }
       timeSpent += 5 * i
       solved++
   }
   fmt.Println(solved)
}
