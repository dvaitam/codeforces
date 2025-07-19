package main

import (
   "fmt"
)

func main() {
   var n, k int
   if _, err := fmt.Scan(&n, &k); err != nil {
       return
   }
   max := (n*n + 1) / 2
   if max < k {
       fmt.Println("NO")
       return
   }
   fmt.Println("YES")
   z := 0
   for i := 0; i < n; i++ {
       row := make([]byte, n)
       for j := 0; j < n; j++ {
           if (z+j)%2 == 0 && k > 0 {
               row[j] = 'L'
               k--
           } else {
               row[j] = 'S'
           }
       }
       z++
       fmt.Println(string(row))
   }
}
