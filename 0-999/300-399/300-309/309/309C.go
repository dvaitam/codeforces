package main

import (
   "fmt"
)

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   // reverse the string
   r := []rune(s)
   for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
       r[i], r[j] = r[j], r[i]
   }
   fmt.Println(string(r))
}
