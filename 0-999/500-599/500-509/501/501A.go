package main

import (
   "fmt"
)

// score computes the Codeforces score for a problem of base points p submitted at minute t.
func score(p, t int) int {
   term1 := 3 * p / 10
   term2 := p - (p/250)*t
   if term1 > term2 {
       return term1
   }
   return term2
}

func main() {
   var a, b, c, d int
   // Read input: a, b are problem base scores; c, d are submission times
   if _, err := fmt.Scan(&a, &b, &c, &d); err != nil {
       return
   }
   s1 := score(a, c)
   s2 := score(b, d)
   switch {
   case s1 > s2:
       fmt.Println("Misha")
   case s2 > s1:
       fmt.Println("Vasya")
   default:
       fmt.Println("Tie")
   }
}
