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

   var n, m, k int
   fmt.Fscan(reader, &n, &m, &k)
   a := make([]int64, n+2)
   for i := 1; i <= n; i++ {
       var v int64
       fmt.Fscan(reader, &v)
       a[i] = v
   }
   l := make([]int, m+2)
   r := make([]int, m+2)
   d := make([]int64, m+2)
   for i := 1; i <= m; i++ {
       fmt.Fscan(reader, &l[i], &r[i], &d[i])
   }
   // Count how many times each operation is applied
   opCount := make([]int64, m+3)
   for i := 0; i < k; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       opCount[x]++
       opCount[y+1]--
   }
   for i := 1; i <= m; i++ {
       opCount[i] += opCount[i-1]
   }
   // Prepare range add for array
   delta := make([]int64, n+3)
   for i := 1; i <= m; i++ {
       if opCount[i] == 0 {
           continue
       }
       add := opCount[i] * d[i]
       delta[l[i]] += add
       delta[r[i]+1] -= add
   }
   // Apply to a
   var curr int64
   for i := 1; i <= n; i++ {
       curr += delta[i]
       a[i] += curr
       // output
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, a[i])
   }
   writer.WriteByte('\n')
}
