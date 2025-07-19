package main

import "fmt"

func main() {
   var n, m int64
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   // minimum value: n - 2*m, not less than 0
   mm := m * 2
   min := n - mm
   if min < 0 {
       min = 0
   }
   // find smallest i such that i*(i-1)/2 >= m
   low, high := int64(0), n+1
   for low < high {
       mid := (low + high) / 2
       if mid*(mid-1)/2 >= m {
           high = mid
       } else {
           low = mid + 1
       }
   }
   // maximum value: n - low
   max := n - low
   fmt.Println(min, max)
}
