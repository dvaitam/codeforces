package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, n)
       for j := 0; j < n; j++ {
           fmt.Scan(&a[i][j])
       }
   }
   // symmetric part: (a[i][j] + a[j][i]) / 2
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           v := float64(a[i][j]+a[j][i]) / 2.0
           if j < n-1 {
               fmt.Printf("%.8f ", v)
           } else {
               fmt.Printf("%.8f", v)
           }
       }
       fmt.Println()
   }
   // skew-symmetric part: (a[i][j] - a[j][i]) / 2
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           v := float64(a[i][j]-a[j][i]) / 2.0
           if j < n-1 {
               fmt.Printf("%.8f ", v)
           } else {
               fmt.Printf("%.8f", v)
           }
       }
       fmt.Println()
   }
}
