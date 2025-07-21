package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   // Map input 1..5 to output 1..3
   // Natural grouping: {1,2}->1, {3,4}->2, {5}->3
   result := (n + 1) / 2
   fmt.Println(result)
}
