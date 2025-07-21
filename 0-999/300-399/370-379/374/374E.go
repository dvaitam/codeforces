package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   blue := make([]struct{u, v int}, n)
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       blue[i].u = x + y
       blue[i].v = x - y
   }
   redU := make([]int, m)
   redV := make([]int, m)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       redU[i] = x + y
       redV[i] = x - y
   }
   // sort blue by u to prune
   sort.Slice(blue, func(i, j int) bool { return blue[i].u < blue[j].u })
   // prepare bitsets
   w := (m + 63) / 64
   A := make([][]uint64, n)
   B := make([][]uint64, n)
   // binary search t
   lo, hi := 0, 2000005
   ans := -1
   for lo <= hi {
       mid := (lo + hi) / 2
       // build bitsets for this mid
       for i := 0; i < n; i++ {
           ai := make([]uint64, w)
           bi := make([]uint64, w)
           ui, vi := blue[i].u, blue[i].v
           // fill bits
           for k := 0; k < m; k++ {
               du := redU[k] - ui
               if du < 0 {
                   du = -du
               }
               dv := redV[k] - vi
               if dv < 0 {
                   dv = -dv
               }
               if du <= 2*mid {
                   ai[k>>6] |= 1 << (uint(k) & 63)
               }
               if dv <= 2*mid {
                   bi[k>>6] |= 1 << (uint(k) & 63)
               }
           }
           A[i] = ai
           B[i] = bi
       }
       ok := false
       // check pairs
       for i := 0; i < n && !ok; i++ {
           ui := blue[i].u
           for j := i + 1; j < n; j++ {
               uj := blue[j].u
               // prune on u distance: need uj-ui <= 4*mid
               if uj - ui > 4*mid {
                   break
               }
               // compute intersection k bits: must satisfy both A and B
               cnt := 0
               for t64 := 0; t64 < w; t64++ {
                   bits := A[i][t64] & A[j][t64] & B[i][t64] & B[j][t64]
                   if bits != 0 {
                       // count bits
                       c := bits
                       // builtin popcount
                       cnt += bitsOn(c)
                       if cnt >= 2 {
                           ok = true
                           break
                       }
                   }
               }
               if ok {
                   break
               }
           }
       }
       if ok {
           ans = mid
           hi = mid - 1
       } else {
           lo = mid + 1
       }
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   if ans < 0 {
       fmt.Fprint(out, "Poor Sereja!")
   } else {
       fmt.Fprint(out, ans)
   }
}

// bitsOn returns number of set bits in x
func bitsOn(x uint64) int {
   // builtin: use Go 1.9+ popcount
   return bits.OnesCount64(x)
}
