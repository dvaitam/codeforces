package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }

   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&a[i])
   }

   var m int
   fmt.Scan(&m)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Scan(&x, &y)
       x--
       left := y - 1
       right := a[x] - y
       a[x] = 0
       if x-1 >= 0 {
           a[x-1] += left
       }
       if x+1 < n {
           a[x+1] += right
       }
   }

   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i := 0; i < n; i++ {
       fmt.Fprintln(w, a[i])
   }
}
