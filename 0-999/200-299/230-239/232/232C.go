package main

import (
   "bufio"
   "fmt"
   "os"
)

func minInt64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func endpointDistances(k int, u uint64, f []uint64, distSE []int64) (int64, int64) {
   // iterative stack for recursion
   ks := make([]int, 0, k+1)
   us := make([]uint64, 0, k+1)
   // descend to base
   for k > 1 {
       ks = append(ks, k)
       us = append(us, u)
       if u <= f[k-1] {
           // left part
           k = k - 1
           // u unchanged
       } else {
           // right part
           u = u - f[k-1]
           k = k - 2
       }
   }
   // base case k <= 1
   var ds, de int64
   if k == 0 {
       ds, de = 0, 0
   } else if k == 1 {
       // D(1): nodes 1 and 2
       if u == 1 {
           ds, de = 0, 1
       } else {
           ds, de = 1, 0
       }
   }
   // unwind
   for i := len(ks) - 1; i >= 0; i-- {
       ck := ks[i]
       cu := us[i]
       // child distances are ds,de from deeper
       if cu <= f[ck-1] {
           // left case: in D(ck-1)
           a, b := ds, de
           // start distance unchanged
           ds = a
           // end via cross to D(ck-2)
           de = minInt64(a, b) + 1 + distSE[ck-2]
       } else {
           // right case: in D(ck-2)
           a2, b2 := ds, de
           ds = a2 + 1
           de = b2
       }
   }
   return ds, de
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t, n int
   fmt.Fscan(reader, &t, &n)
   // precompute Fibonacci-like sizes
   maxU := uint64(10000000000000000) // 1e16
   f := make([]uint64, n+1)
   f[0], f[1] = 1, 2
   for i := 2; i <= n; i++ {
       v := f[i-1] + f[i-2]
       if v > maxU {
           f[i] = maxU + 1
       } else {
           f[i] = v
       }
   }
   // precompute distSE: distance from start to end in D(k)
   distSE := make([]int64, n+1)
   if n >= 0 {
       distSE[0] = 0
   }
   if n >= 1 {
       distSE[1] = 1
   }
   for i := 2; i <= n; i++ {
       // option A: via start cross
       a := 1 + distSE[i-2]
       // option B: via end cross
       b := distSE[i-1] + 1 + distSE[i-2]
       if a < b {
           distSE[i] = a
       } else {
           distSE[i] = b
       }
   }
   // process queries
   for qi := 0; qi < t; qi++ {
       var u, v uint64
       fmt.Fscan(reader, &u, &v)
       if u == v {
           fmt.Fprintln(writer, 0)
           continue
       }
       // copy n
       kk := n
       var ans int64
       // iterative reduction
       for {
           if kk == 1 {
               ans = 1
               break
           }
           if u <= f[kk-1] && v <= f[kk-1] {
               kk--
               continue
           }
           if u > f[kk-1] && v > f[kk-1] {
               u -= f[kk-1]
               v -= f[kk-1]
               kk -= 2
               continue
           }
           // cross-case
           if u > f[kk-1] {
               u, v = v, u
           }
           // u in left, v in right
           // distances for u in D(kk-1)
           du1, du2 := endpointDistances(kk-1, u, f, distSE)
           // distances for v' in D(kk-2)
           v2 := v - f[kk-1]
           dv1, _ := endpointDistances(kk-2, v2, f, distSE)
           ans = minInt64(minInt64(du1, du2)+1+dv1, 1<<60)
           break
       }
       fmt.Fprintln(writer, ans)
   }
}
