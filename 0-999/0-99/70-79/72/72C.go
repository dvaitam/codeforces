package main

import "fmt"

func main() {
   var x int
   if _, err := fmt.Scan(&x); err != nil {
       return
   }
   // x is extraordinarily nice if it has exactly one factor 2 (i.e., divisible by 2 but not by 4)
   if x%2 == 0 && x%4 != 0 {
       fmt.Println("yes")
   } else {
       fmt.Println("no")
   }
}
