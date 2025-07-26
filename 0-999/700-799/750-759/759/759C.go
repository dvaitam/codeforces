package main

import "fmt"

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   // Reverse the string by runes
   rs := []rune(s)
   for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
       rs[i], rs[j] = rs[j], rs[i]
   }
   fmt.Println(string(rs))
}
