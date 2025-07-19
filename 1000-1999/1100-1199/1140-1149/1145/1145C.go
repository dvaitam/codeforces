package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   arr := []int{14, 12, 13, 8, 9, 10, 11, 0, 1, 2, 3, 4, 5, 6, 7}
   if n >= 1 && n <= len(arr) {
       fmt.Println(arr[n-1])
   } else {
       fmt.Println(15)
   }
}
