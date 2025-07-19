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
   A := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &A[i])
   }
   half := (n + 1) / 2
   for i := -1000; i < 1000; i++ {
       if i == 0 {
           continue
       }
       cnt := 0
       for _, v := range A {
           if v*i > 0 {
               cnt++
           }
       }
       if cnt >= half {
           fmt.Println(i)
           return
       }
   }
   fmt.Println(0)
}
