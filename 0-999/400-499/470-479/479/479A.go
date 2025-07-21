package main

import "fmt"

// max returns the maximum value among the provided integers.
func max(nums ...int) int {
   m := nums[0]
   for _, v := range nums {
       if v > m {
           m = v
       }
   }
   return m
}

func main() {
   var a, b, c int
   // Read three integers from standard input
   if _, err := fmt.Scan(&a, &b, &c); err != nil {
       return
   }

   // Evaluate all possible expressions by inserting '+' or '*' and parentheses
   results := []int{
       a + b + c,
       a * b * c,
       (a + b) * c,
       a * (b + c),
       a * b + c,
       a + b * c,
   }

   // Print the maximum result
   fmt.Println(max(results...))
}
