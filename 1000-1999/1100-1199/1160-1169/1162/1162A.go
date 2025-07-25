package main

import (
   "fmt"
)

func main() {
   var n, h, m int
   if _, err := fmt.Scan(&n, &h, &m); err != nil {
       return
   }
   heights := make([]int, n)
   for i := range heights {
       heights[i] = h
   }
   for i := 0; i < m; i++ {
       var l, r, x int
       fmt.Scan(&l, &r, &x)
       for j := l - 1; j < r; j++ {
           if heights[j] > x {
               heights[j] = x
           }
       }
   }
   total := 0
   for _, v := range heights {
       total += v * v
   }
   fmt.Println(total)
