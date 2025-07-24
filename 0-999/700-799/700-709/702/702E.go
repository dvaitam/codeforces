package main

import (
   "bufio"
   "fmt"
   "math"
   "math/bits"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   var k uint64
   fmt.Fscan(reader, &n, &k)
   f := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &f[i])
   }
   // initial weights
   sum0 := make([]int64, n)
   min0 := make([]int64, n)
   for i := 0; i < n; i++ {
       var w int64
       fmt.Fscan(reader, &w)
       sum0[i] = w
       min0[i] = w
   }
   // number of bits needed for k
   maxJ := bits.Len64(k)
   // parent[j][i]: node reached from i by 2^j steps
   parent := make([][]int, maxJ)
   sum := make([][]int64, maxJ)
   minw := make([][]int64, maxJ)
   // level 0 initialization
   parent[0] = make([]int, n)
   sum[0] = make([]int64, n)
   minw[0] = make([]int64, n)
   for i := 0; i < n; i++ {
       parent[0][i] = f[i]
       sum[0][i] = sum0[i]
       minw[0][i] = min0[i]
   }
   // binary lifting precompute
   for j := 1; j < maxJ; j++ {
       parent[j] = make([]int, n)
       sum[j] = make([]int64, n)
       minw[j] = make([]int64, n)
       for i := 0; i < n; i++ {
           pi := parent[j-1][i]
           parent[j][i] = parent[j-1][pi]
           sum[j][i] = sum[j-1][i] + sum[j-1][pi]
           // minimal weight on this jump
           m1 := minw[j-1][i]
           m2 := minw[j-1][pi]
           if m2 < m1 {
               m1 = m2
           }
           minw[j][i] = m1
       }
   }
   // answer queries for each start
   for i := 0; i < n; i++ {
       cur := i
       var total int64
       bestMin := int64(math.MaxInt64)
       kk := k
       j := 0
       for kk > 0 {
           if kk&1 == 1 {
               total += sum[j][cur]
               if minw[j][cur] < bestMin {
                   bestMin = minw[j][cur]
               }
               cur = parent[j][cur]
           }
           kk >>= 1
           j++
       }
       fmt.Fprintln(writer, total, bestMin)
   }
}
