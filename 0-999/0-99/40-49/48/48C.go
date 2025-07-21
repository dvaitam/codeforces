package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   stops := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &stops[i])
   }
   // Determine valid range for alpha: intersection of [L_i, R_i)
   L := 10.0 * float64(stops[0]) / 1.0
   R := (10.0*float64(stops[0]) + 10.0) / 1.0
   for i := 2; i <= n; i++ {
       si := float64(stops[i-1])
       li := 10.0 * si / float64(i)
       ri := (10.0*si + 10.0) / float64(i)
       if li > L {
           L = li
       }
       if ri < R {
           R = ri
       }
   }
   // Ensure at least alpha >= 10
   if L < 10.0 {
       L = 10.0
   }
   // Compute possible next stop range
   last := stops[n-1]
   // u = (alpha*(n+1) - 10) / 10
   // for alpha in [L, R), u in [u_min, u_max)
   eps := 1e-12
   uMin := ((float64(n+1) * L) - 10.0) / 10.0
   uMax := ((float64(n+1) * R) - 10.0) / 10.0
   // minimal next stop
   mMin := int(floor(uMin + eps))
   cand1 := mMin + 1
   if cand1 < last+1 {
       cand1 = last + 1
   }
   // maximal next stop
   mMax1 := int(ceil(uMax - eps))
   cand2 := mMax1
   if cand2 < last+1 {
       cand2 = last + 1
   }
   // Output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   if cand1 == cand2 {
       fmt.Fprintln(writer, "unique")
       fmt.Fprintln(writer, cand1)
   } else {
       fmt.Fprintln(writer, "not unique")
   }
}

func floor(x float64) float64 {
   // simple wrapper
   return float64(int(x)) - func() float64 {
       if x < 0 && float64(int(x)) != x {
           return 1
       }
       return 0
   }()
}

func ceil(x float64) float64 {
   // simple wrapper
   return float64(int(x)) + func() float64 {
       if x > float64(int(x)) {
           return 1
       }
       return 0
   }()
}
