package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   var s string
   fmt.Scan(&s)
   // Use a slice of runes to build the original word
   res := make([]rune, 0, n)
   runes := []rune(s)
   // Reverse the encoding: insert characters back in reverse order
   for i := n - 1; i >= 0; i-- {
       idx := len(res) / 2
       // insert runes[i] at position idx
       res = append(res[:idx], append([]rune{runes[i]}, res[idx:]...)...)
   }
   fmt.Println(string(res))
}
