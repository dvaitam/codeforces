package main

import (
   "fmt"
)

func main() {
   var m, n int
   if _, err := fmt.Scan(&m, &n); err != nil {
       return
   }
   B := make([][]int, m)
   for i := 0; i < m; i++ {
       B[i] = make([]int, n)
       for j := 0; j < n; j++ {
           fmt.Scan(&B[i][j])
       }
   }
   A := make([][]int, m)
   for i := 0; i < m; i++ {
       A[i] = make([]int, n)
       for j := 0; j < n; j++ {
           okRow := true
           for k := 0; k < n; k++ {
               if B[i][k] == 0 {
                   okRow = false
                   break
               }
           }
           if !okRow {
               continue
           }
           okCol := true
           for k := 0; k < m; k++ {
               if B[k][j] == 0 {
                   okCol = false
                   break
               }
           }
           if okCol {
               A[i][j] = 1
           }
       }
   }
   // Reconstruct B from A
   Brest := make([][]int, m)
   for i := 0; i < m; i++ {
       Brest[i] = make([]int, n)
   }
   for i := 0; i < m; i++ {
       for j := 0; j < n; j++ {
           if A[i][j] == 1 {
               for k := 0; k < n; k++ {
                   Brest[i][k] = 1
               }
               for k := 0; k < m; k++ {
                   Brest[k][j] = 1
               }
           }
       }
   }
   // Verify
   for i := 0; i < m; i++ {
       for j := 0; j < n; j++ {
           if Brest[i][j] != B[i][j] {
               fmt.Println("NO")
               return
           }
       }
   }
   fmt.Println("YES")
   for i := 0; i < m; i++ {
       for j := 0; j < n; j++ {
           if j > 0 {
               fmt.Print(" ")
           }
           fmt.Print(A[i][j])
       }
       fmt.Println()
   }
}
