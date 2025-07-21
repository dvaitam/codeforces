package main

import (
   "fmt"
   "sort"
)

func sumDigits(x int64) int {
   sum := 0
   for x > 0 {
       sum += int(x % 10)
       x /= 10
   }
   return sum
}

func main() {
   var a, b, c int
   if _, err := fmt.Scan(&a, &b, &c); err != nil {
       return
   }
   var res []int
   // sum of digits s ranges from 1 to 81 (max for 9-digit numbers)
   for s := 1; s <= 81; s++ {
       // compute s^a
       p := int64(1)
       s64 := int64(s)
       for i := 0; i < a; i++ {
           p *= s64
       }
       x64 := int64(b)*p + int64(c)
       if x64 <= 0 || x64 >= 1000000000 {
           continue
       }
       if sumDigits(x64) == s {
           res = append(res, int(x64))
       }
   }
   sort.Ints(res)
   fmt.Println(len(res))
   if len(res) > 0 {
       for i, v := range res {
           if i > 0 {
               fmt.Print(" ")
           }
           fmt.Print(v)
       }
       fmt.Println()
   }
}
