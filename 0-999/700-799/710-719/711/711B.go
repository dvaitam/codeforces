package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([][]int64, n)
   var zi, zj int
   for i := 0; i < n; i++ {
       a[i] = make([]int64, n)
       for j := 0; j < n; j++ {
           fmt.Fscan(in, &a[i][j])
           if a[i][j] == 0 {
               zi, zj = i, j
           }
       }
   }
   // Special case: 1x1 grid
   if n == 1 {
       fmt.Println(1)
       return
   }
   // Determine target sum S from a complete row
   var S int64
   found := false
   for i := 0; i < n; i++ {
       if i == zi {
           continue
       }
       var sum int64
       for j := 0; j < n; j++ {
           sum += a[i][j]
       }
       if !found {
           S = sum
           found = true
       } else if sum != S {
           fmt.Println(-1)
           return
       }
   }
   // Compute the required value for the empty cell
   var sumZi int64
   for j := 0; j < n; j++ {
       sumZi += a[zi][j]
   }
   x := S - sumZi
   if x <= 0 {
       fmt.Println(-1)
       return
   }
   a[zi][zj] = x
   // Verify all rows
   for i := 0; i < n; i++ {
       var sum int64
       for j := 0; j < n; j++ {
           sum += a[i][j]
       }
       if sum != S {
           fmt.Println(-1)
           return
       }
   }
   // Verify all columns
   for j := 0; j < n; j++ {
       var sum int64
       for i := 0; i < n; i++ {
           sum += a[i][j]
       }
       if sum != S {
           fmt.Println(-1)
           return
       }
   }
   // Verify main diagonal
   var diagSum int64
   for i := 0; i < n; i++ {
       diagSum += a[i][i]
   }
   if diagSum != S {
       fmt.Println(-1)
       return
   }
   // Verify secondary diagonal
   diagSum = 0
   for i := 0; i < n; i++ {
       diagSum += a[i][n-1-i]
   }
   if diagSum != S {
       fmt.Println(-1)
       return
   }
   fmt.Println(x)
}
