package main

import (
   "fmt"
   "sort"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   arr := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&arr[i])
   }
   sort.Ints(arr)
   for i, v := range arr {
       if i > 0 {
           fmt.Print(" ")
       }
       fmt.Print(v)
   }
   fmt.Println()
}
