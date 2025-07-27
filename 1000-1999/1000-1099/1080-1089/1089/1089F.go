package main

import "fmt"

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   result := 1
   for i := 2; i <= n; i++ {
       result *= i
   }
   fmt.Println(result)
}
