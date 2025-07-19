package main

import (
   "fmt"
)

func main() {
   var N, A, B int
   if _, err := fmt.Scan(&N, &A, &B); err != nil {
       return
   }
   if N > A*B {
       fmt.Println(-1)
       return
   }
   // initialize matrix
   M := make([][]int, A)
   for i := 0; i < A; i++ {
       M[i] = make([]int, B)
   }
   cnt := 1
   // fill matrix with pattern
   for i := 0; i < A; i++ {
       for j := 0; j < B; j++ {
           M[i][j] = cnt
           if B%2 == 0 && i%2 != 0 {
               if j%2 == 0 {
                   M[i][j]++
               } else {
                   M[i][j]--
               }
           }
           if M[i][j] > N {
               M[i][j] = 0
           }
           cnt++
           if cnt > N+1 {
               break
           }
       }
       if cnt > N+1 {
           break
       }
   }
   // output matrix
   for i := 0; i < A; i++ {
       for j := 0; j < B; j++ {
           if j == B-1 {
               fmt.Printf("%d\n", M[i][j])
           } else {
               fmt.Printf("%d ", M[i][j])
           }
       }
   }
}
