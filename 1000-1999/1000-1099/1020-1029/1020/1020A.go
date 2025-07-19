package main

import "fmt"

// abs returns the absolute value of x.
func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   var n, h, a, b, q int
   // Read building count n, height h, range [a,b], number of queries q
   if _, err := fmt.Scan(&n, &h, &a, &b, &q); err != nil {
       return
   }
   for i := 0; i < q; i++ {
       var n1, h1, n2, h2 int
       fmt.Scan(&n1, &h1, &n2, &h2)
       if n1 == n2 {
           // Same building: only vertical distance
           fmt.Println(abs(h2 - h1))
       } else {
           // Move to nearest valid floor [a,b] in building n1
           bsv := h1
           if h1 > b {
               bsv = b
           } else if h1 < a {
               bsv = a
           }
           // Total distance: vertical move in n1, horizontal move, vertical move in n2
           ans := abs(h2 - bsv) + abs(bsv - h1) + abs(n2 - n1)
           fmt.Println(ans)
       }
   }
}
