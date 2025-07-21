package main

import "fmt"

// distinct checks if all digits in n are unique
func distinct(n int) bool {
   var seen [10]bool
   for n > 0 {
       d := n % 10
       if seen[d] {
           return false
       }
       seen[d] = true
       n /= 10
   }
   return true
}

func main() {
   var y int
   if _, err := fmt.Scan(&y); err != nil {
       return
   }
   for i := y + 1; ; i++ {
       if distinct(i) {
           fmt.Println(i)
           break
       }
   }
}
