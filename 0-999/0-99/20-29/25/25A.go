package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       if _, err := fmt.Scan(&a[i]); err != nil {
           return
       }
   }
   cntEven := 0
   for _, v := range a {
       if v%2 == 0 {
           cntEven++
       }
   }
   wantEven := cntEven == 1
   for i, v := range a {
       if (v%2 == 0) == wantEven {
           fmt.Println(i + 1)
           return
       }
   }
}
