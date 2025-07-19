package main

import "fmt"

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   a := make([]int, n)
   b := make([]int, n)
   c := make([]int, n)
   d := make([]int, n)
   vis := make([]bool, n)
   for i := 0; i < n; i++ {
       if _, err := fmt.Scan(&a[i], &b[i], &c[i], &d[i]); err != nil {
           return
       }
   }
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if a[i] < a[j] && b[i] < b[j] && c[i] < c[j] {
               vis[i] = true
           }
       }
   }
   ans := 1005
   idx := 0
   for i := 0; i < n; i++ {
       if !vis[i] && d[i] < ans {
           ans = d[i]
           idx = i
       }
   }
   fmt.Println(idx + 1)
}
