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

   var T int
   fmt.Fscan(reader, &T)
   const INF = int64(1e18)

   for T > 0 {
       T--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       lim := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &lim[i])
       }
       // options per task, include dummy at index 0
       a := make([][][3]int64, n)
       for i := 0; i < n; i++ {
           a[i] = make([][3]int64, 1)
       }
       for id := int64(0); id < int64(m); id++ {
           var e, t, p int64
           fmt.Fscan(reader, &e, &t, &p)
           e--
           a[e] = append(a[e], [3]int64{t, p, id})
       }
       var ans []int64
       var cur int64
       ok := true
       // process each task
       for i := 0; i < n; i++ {
           opts := a[i]
           N := len(opts) - 1
           // dp[j][k]: min time to achieve k percent using first j options
           dp := make([][]int64, N+1)
           for j := 0; j <= N; j++ {
               dp[j] = make([]int64, 101)
               for k := 1; k <= 100; k++ {
                   dp[j][k] = INF
               }
               dp[j][0] = 0
           }
           for j := 1; j <= N; j++ {
               t := opts[j][0]
               p := opts[j][1]
               for k := 0; k <= 100; k++ {
                   // not take
                   v := dp[j-1][k]
                   // take
                   kp := k - int(p)
                   if kp < 0 {
                       kp = 0
                   }
                   if dp[j-1][kp]+t < v {
                       v = dp[j-1][kp] + t
                   }
                   dp[j][k] = v
               }
           }
           need := dp[N][100]
           cur += need
           if cur > lim[i] {
               fmt.Fprintln(writer, -1)
               ok = false
               break
           }
           // backtrack
           k := 100
           for j := N; j >= 1; j-- {
               t := opts[j][0]
               p := opts[j][1]
               id := opts[j][2]
               kp := k - int(p)
               if kp < 0 {
                   kp = 0
               }
               if dp[j][k] == dp[j-1][kp]+t {
                   ans = append(ans, id)
                   k = kp
               }
           }
       }
       if !ok {
           continue
       }
       // output
       fmt.Fprintln(writer, len(ans))
       for i, id := range ans {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, id+1)
       }
       fmt.Fprintln(writer)
   }
}
