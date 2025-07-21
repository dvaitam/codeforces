package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   pos := 1
   // make n-1 throws, distances 1,2,...,n-1
   for i := 1; i < n; i++ {
       pos = (pos + i - 1) % n + 1
       fmt.Print(pos)
       if i != n-1 {
           fmt.Print(" ")
       }
   }
   // newline after output
   fmt.Println()
}
