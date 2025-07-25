package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for i := 0; i < t; i++ {
       var n int64
       var k int64
       fmt.Fscan(in, &n, &k)
       if k == 0 {
           // no splits, keep whole square
           fmt.Fprintf(out, "YES %d\n", n)
           continue
       }
       // precompute Mmin for d from 0 to maxD
       // find maximum d>=1 such that Mmin = (4^d -1)/3 <= k
       maxD := int64(0)
       // since for d>=32, Mmin >1e18, limit to 60 but break early
       for d := int64(1); d <= 60 && d <= n; d++ {
           // compute 4^d
           var pow4_d int64 = 1
           for j := int64(0); j < d; j++ {
               pow4_d *= 4
               if pow4_d < 0 {
                   break
               }
           }
           // Mmin = (pow4_d -1)/3
           mmin := (pow4_d - 1) / 3
           if mmin <= k {
               maxD = d
           } else {
               break
           }
       }
       found := false
       var ansS int64
       // try d from maxD down to 1
       for d := maxD; d >= 1; d-- {
           // compute Mmin
           var pow4_d int64 = 1
           for j := int64(0); j < d; j++ {
               pow4_d *= 4
           }
           mmin := (pow4_d - 1) / 3
           if mmin > k {
               continue
           }
           // check Tmax
           // for d <= n-31, Tmax >= k automatically (and d>=1)
           if n-d > 30 {
               // enough large
               ansS = n - d
               found = true
               break
           }
           // compute C = 2^d
           var C int64 = 1
           for j := int64(0); j < d; j++ {
               C <<= 1
           }
           // B = 4^(n-d)
           s := n - d
           var B int64 = 1
           for j := int64(0); j < s; j++ {
               B *= 4
           }
           // Tmax = (B*(C-1)^2 + 2*(C-1)) / 3
           cm1 := C - 1
           // use int128 via Go big int not needed: values fit in int64 for this branch
           tnum := B*cm1*cm1 + 2*cm1
           tmax := tnum / 3
           if k <= tmax {
               ansS = n - d
               found = true
               break
           }
       }
       if found {
           fmt.Fprintf(out, "YES %d\n", ansS)
       } else {
           fmt.Fprintln(out, "NO")
       }
   }
}
