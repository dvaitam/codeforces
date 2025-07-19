package main

import (
   "fmt"
)

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   inSet := map[rune]bool{
       'A': true, 'E': true, 'F': true, 'H': true,
       'I': true, 'K': true, 'L': true, 'M': true,
       'N': true, 'T': true, 'V': true, 'W': true,
       'X': true, 'Y': true, 'Z': true,
   }
   allIn, noneIn := true, true
   for _, c := range s {
       if inSet[c] {
           noneIn = false
       } else {
           allIn = false
       }
   }
   if allIn || noneIn {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
