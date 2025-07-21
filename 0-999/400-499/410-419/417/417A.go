package main

import (
   "fmt"
)

func main() {
   var c, d, n, m, k int
   if _, err := fmt.Scan(&c, &d); err != nil {
       return
   }
   fmt.Scan(&n, &m)
   fmt.Scan(&k)
   totalNeeded := n*m - k
   if totalNeeded < 0 {
       totalNeeded = 0
   }
   // maximum main rounds needed to cover all
   maxMain := 0
   if n > 0 {
       maxMain = (totalNeeded + n - 1) / n
   }
   minProblems := int(^uint(0) >> 1) // max int
   for x := 0; x <= maxMain; x++ {
       rem := totalNeeded - x*n
       if rem < 0 {
           rem = 0
       }
       cost := x*c + rem*d
       if cost < minProblems {
           minProblems = cost
       }
   }
   fmt.Println(minProblems)
}
