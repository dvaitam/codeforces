package main

import "fmt"

func main() {
   var v int
   if _, err := fmt.Scan(&v); err != nil {
       return
   }
   // Build a list of distinct integers 1,2,... until sum exceeds v
   a := make([]int, 0, v)
   sum := 0
   n := 0
   for {
       n++
       a = append(a, n)
       sum += n
       if sum > v {
           break
       }
   }
   // Remove the element to adjust the sum exactly to v
   diffIndex := sum - v - 1
   if diffIndex >= 0 && diffIndex < len(a) {
       a = append(a[:diffIndex], a[diffIndex+1:]...)
   }
   // Output the results
   fmt.Println(len(a))
   for i, x := range a {
       if i > 0 {
           fmt.Print(" ")
       }
       fmt.Print(x)
   }
   if len(a) > 0 {
       fmt.Println()
   }
}
