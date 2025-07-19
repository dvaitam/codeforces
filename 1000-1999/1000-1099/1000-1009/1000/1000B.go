package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var m int64
   fmt.Fscan(reader, &n, &m)
   // a[0]=0, a[1..n]=input, a[n+1]=m
   a := make([]int64, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   a[n+1] = m
   N := n + 1
   // sum[i]: total lit time up to a[i]
   sum := make([]int64, N+1)
   for i := 1; i <= N; i++ {
       if i%2 == 1 {
           sum[i] = sum[i-1] + (a[i] - a[i-1])
       } else {
           sum[i] = sum[i-1]
       }
   }
   ans := sum[N]
   for i := 1; i < N; i++ {
       if i%2 == 1 {
           // try inserting between a[i] and a[i+1]
           if a[i]+1 < a[i+1] {
               ret := sum[i] + (a[i+1] - a[i] - 1)
               tmp := m - a[i+1] - sum[N] + sum[i+1]
               ret += tmp
               if ret > ans {
                   ans = ret
               }
           }
           // try inserting before a[i]
           if a[i]-1 > a[i-1] {
               ret := sum[i] - 1
               tmp := m - a[i] - sum[N] + sum[i]
               ret += tmp
               if ret > ans {
                   ans = ret
               }
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   fmt.Fprint(writer, ans)
   writer.Flush()
}
