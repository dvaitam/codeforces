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

   var t int
   fmt.Fscan(reader, &t)
   const INF int64 = 1000000000000000000
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       sum := make([]int64, n+1)
       for i := 0; i < n; i++ {
           sum[i+1] = sum[i] + a[i]
       }
       mins := make([]int64, n+1)
       copy(mins, sum)
       for i := n - 1; i >= 0; i-- {
           if mins[i+1] < mins[i] {
               mins[i] = mins[i+1]
           }
       }
       maxv := sum[n]
       ans := INF
       for i := 0; i < n; i++ {
           diff := sum[i] - mins[i]
           if diff < 0 {
               diff = 0
           }
           v := sum[n] + diff
           if v > maxv {
               maxv = v
               ans = sum[i]
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
