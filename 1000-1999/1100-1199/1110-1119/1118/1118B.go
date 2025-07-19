package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sum := make([]int64, n+2)
   for i := 1; i <= n; i++ {
       sum[i] = sum[i-1] + a[i]
   }
   s0 := make([]int64, n+2)
   for i := 2; i <= n; i++ {
       s0[i] = sum[i-1] - s0[i-1]
   }
   s1 := make([]int64, n+3)
   for i := n - 1; i >= 1; i-- {
       s1[i] = sum[n] - sum[i] - s1[i+1]
   }
   ans := 0
   for i := 1; i <= n; i++ {
       if s0[i] + s1[i+1] == s1[i] + s0[i-1] {
           ans++
       }
   }
   fmt.Fprint(writer, ans)
}
