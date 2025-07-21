package main

import (
   "fmt"
)

func canHold(girl, boy int) bool {
   // No two girl's fingers touch: need boy >= girl-1
   // No three boy's fingers touch: need boy <= 2*(girl+1)
   return boy >= girl-1 && boy <= 2*(girl+1)
}

func main() {
   var aL, aR, bL, bR int
   if _, err := fmt.Scan(&aL, &aR); err != nil {
       return
   }
   if _, err := fmt.Scan(&bL, &bR); err != nil {
       return
   }
   // Check pairing: girl left with boy right, or girl right with boy left
   if canHold(aL, bR) || canHold(aR, bL) {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
