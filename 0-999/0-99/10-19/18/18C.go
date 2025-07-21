package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var total int64
   for _, v := range a {
       total += v
   }
   var prefix int64
   var count int64
   for i := 0; i < n-1; i++ {
       prefix += a[i]
       if prefix*2 == total {
           count++
       }
   }
   fmt.Println(count)
}
