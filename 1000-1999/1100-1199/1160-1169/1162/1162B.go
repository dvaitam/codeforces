package main

import (
   "fmt"
)

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   a := make([][]int, n)
   b := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Scan(&a[i][j])
       }
   }
   for i := 0; i < n; i++ {
       b[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Scan(&b[i][j])
       }
   }
   // Swap to ensure a[i][j] >= b[i][j]
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if a[i][j] < b[i][j] {
               a[i][j], b[i][j] = b[i][j], a[i][j]
           }
       }
   }
   if isIncreasing(a, n, m) && isIncreasing(b, n, m) {
       fmt.Println("Possible")
   } else {
       fmt.Println("Impossible")
   }
}

// isIncreasing checks that each row and column is strictly increasing
func isIncreasing(mat [][]int, n, m int) bool {
   // Rows
   for i := 0; i < n; i++ {
       for j := 1; j < m; j++ {
           if mat[i][j] <= mat[i][j-1] {
               return false
           }
       }
   }
   // Columns
   for j := 0; j < m; j++ {
       for i := 1; i < n; i++ {
           if mat[i][j] <= mat[i-1][j] {
               return false
           }
       }
   }
   return true
}
