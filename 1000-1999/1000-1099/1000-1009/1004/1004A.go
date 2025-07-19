package main

import (
   "fmt"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   var n, d int
   if _, err := fmt.Scan(&n, &d); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&a[i])
   }
   t := 0
   for i := 0; i < n; i++ {
       if i != 0 {
           if abs(a[i]-d - a[i-1]) >= d && a[i]-a[i-1] > 2*d {
               t++
           }
       }
       if i != n-1 {
           if abs(a[i]+d - a[i+1]) >= d {
               t++
           }
       }
   }
   fmt.Println(t + 2)
}
