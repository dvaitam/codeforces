package main

import "fmt"

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   result := "EASY"
   for i := 0; i < n; i++ {
       var x int
       if _, err := fmt.Scan(&x); err != nil {
           return
       }
       if x == 1 {
           result = "HARD"
           break
       }
   }
   fmt.Println(result)
}
