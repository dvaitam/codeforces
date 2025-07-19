package main

import (
   "fmt"
   "sort"
)

func main() {
   var nums [4]int
   for i := 0; i < 4; i++ {
       _, err := fmt.Scan(&nums[i])
       if err != nil {
           return
       }
   }
   s := nums[:]
   sort.Ints(s)
   max := s[3]
   // output differences between max and other three numbers
   fmt.Println(max-s[0], max-s[1], max-s[2])
}
