package main

import (
   "fmt"
)

// reverse returns the integer obtained by reversing the digits of n.
func reverse(n int64) int64 {
   var rev int64
   for n > 0 {
       rev = rev*10 + n%10
       n /= 10
   }
   return rev
}

func main() {
   var a, b int64
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   fmt.Println(a + reverse(b))
}
