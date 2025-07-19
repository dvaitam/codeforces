package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var k int64
   fmt.Fscan(reader, &n, &k)
   x := make([]int64, n)
   y := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &y[i])
   }
   var u, v int64
   ok := true
   for i := 0; i < n; i++ {
       a, b := x[i], y[i]
       if (b+1)*k-u < a || (a+1)*k-v < b {
           ok = false
           break
       }
       if b*k-u < a {
           u = a - b*k + u
       } else {
           u = 0
       }
       if a*k-v < b {
           v = b - a*k + v
       } else {
           v = 0
       }
   }
   if ok {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
