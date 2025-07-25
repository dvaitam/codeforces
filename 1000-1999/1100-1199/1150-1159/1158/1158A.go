package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   b := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   g := make([]int64, m)
   for j := 0; j < m; j++ {
       fmt.Fscan(reader, &g[j])
   }

   sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
   sort.Slice(g, func(i, j int) bool { return g[i] < g[j] })

   bMax := b[n-1]
   // if any girl's max less than bMax, impossible
   if g[0] < bMax {
       fmt.Fprintln(writer, -1)
       return
   }

   var sumB int64
   for i := 0; i < n; i++ {
       sumB += b[i]
   }
   // base: each boy gives at least b[i] to each of m girls
   ans := sumB * int64(m)
   // for each girl, raise from bMax to g[j]
   for j := 0; j < m; j++ {
       ans += g[j] - bMax
   }
   // if smallest g > bMax, need one extra adjustment from second max
   if g[0] > bMax {
       // second largest b
       bSec := b[n-2]
       ans += bMax - bSec
   }
   fmt.Fprintln(writer, ans)
}
