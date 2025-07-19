package main

import (
   "fmt"
)

func main() {
   var a uint64
   if _, err := fmt.Scan(&a); err != nil {
       return
   }
   count := 0
   for a > 0 {
       d := a % 10
       if d == 4 || d == 7 {
           count++
       }
       a /= 10
   }
   if count == 0 {
       fmt.Println("NO")
       return
   }
   for tmp := count; tmp > 0; tmp /= 10 {
       d := tmp % 10
       if d != 4 && d != 7 {
           fmt.Println("NO")
           return
       }
   }
   fmt.Println("YES")
}
