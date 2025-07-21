package main

import (
   "fmt"
)

// builds a tournament graph of n vertices with diameter at most 2
func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   if n < 3 {
       fmt.Println(-1)
       return
   }
   // circulant tournament: edge i->j iff (j-i mod n) in [1..n/2]
   a := make([][]int, n)
   for i := range a {
       a[i] = make([]int, n)
       for j := range a[i] {
           if i == j {
               a[i][j] = 0
               continue
           }
           d := (j - i + n) % n
           if d >= 1 && d <= n/2 {
               a[i][j] = 1
           }
       }
   }
   // output adjacency matrix
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if j > 0 {
               fmt.Print(" ")
           }
           fmt.Print(a[i][j])
       }
       fmt.Println()
   }
}
