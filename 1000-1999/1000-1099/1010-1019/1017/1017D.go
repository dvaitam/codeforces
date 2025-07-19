package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, q int
   fmt.Fscan(reader, &n, &m, &q)
   w := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &w[i])
   }
   u := (1 << n) - 1
   sw := make([]int, u+1)
   // sum of weights for each mask
   for mask := 1; mask <= u; mask++ {
       low := mask & -mask
       bit := bits.TrailingZeros(uint(low))
       sw[mask] = sw[mask^low] + w[bit]
   }
   bas := make([]int, u+1)
   // read m bitstrings
   for i := 0; i < m; i++ {
       var s string
       fmt.Fscan(reader, &s)
       mask := 0
       for j, c := range s {
           if c == '1' {
               mask |= 1 << j
           }
       }
       bas[mask]++
   }
   // precompute counts
   cnt := make([][]int, u+1)
   for i := range cnt {
       cnt[i] = make([]int, 101)
   }
   for i := 0; i <= u; i++ {
       for j := 0; j <= u; j++ {
           v := i ^ j ^ u
           if sw[v] <= 100 {
               cnt[i][sw[v]] += bas[j]
           }
       }
       for k := 1; k <= 100; k++ {
           cnt[i][k] += cnt[i][k-1]
       }
   }
   // answer queries
   for i := 0; i < q; i++ {
       var s string
       var k int
       fmt.Fscan(reader, &s, &k)
       mask := 0
       for j, c := range s {
           if c == '1' {
               mask |= 1 << j
           }
       }
       res := cnt[mask][k]
       fmt.Fprintln(writer, res)
   }
}
