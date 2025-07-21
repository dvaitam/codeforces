package main

import "fmt"

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   if n >= 0 {
       fmt.Println(n)
       return
   }
   // Option 1: remove last digit
   a := n / 10
   // Option 2: remove second last digit
   b := (n/100) * 10 + n%10
   if a > b {
       fmt.Println(a)
   } else {
       fmt.Println(b)
   }
}
