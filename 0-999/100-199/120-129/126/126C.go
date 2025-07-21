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
   a := make([][]byte, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(in, &s)
       row := make([]byte, n)
       for j := 0; j < n; j++ {
           // '0'->0, '1'->1
           row[j] = s[j] - '0'
       }
       a[i] = row
   }
   offsum := 0
   // total ones in off-diagonal rows and columns
   totalRow := make([]int, n)
   totalCol := make([]int, n)
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if i == j {
               continue
           }
           if a[i][j] != 0 {
               offsum++
               totalRow[i]++
               totalCol[j]++
           }
       }
   }
   diagops := 0
   for i := 0; i < n; i++ {
       // compute needed diag operation t[i][i]
       t := (int(a[i][i]) + totalRow[i] + totalCol[i]) & 1
       diagops += t
   }
   ans := offsum + diagops
   fmt.Println(ans)
}
