package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   rowSum := make([]int, n)
   colSum := make([]int, n)
   // Read board and compute row and column sums
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           var x int
           fmt.Scan(&x)
           rowSum[i] += x
           colSum[j] += x
       }
   }
   // Count winning squares: colSum[j] > rowSum[i]
   cnt := 0
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if colSum[j] > rowSum[i] {
               cnt++
           }
       }
   }
   fmt.Println(cnt)
}
