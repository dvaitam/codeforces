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
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // binary search max k such that minimal sum_j <= a[j]
   maxK := m / n
   lo, hi := 0, maxK+1
   // check function
   check := func(k int) bool {
       if k == 0 {
           return true
       }
       // base term: n * k * (k-1) / 2
       kk := int64(k)
       base := int64(n) * kk * (kk - 1) / 2
       for j := 0; j < n; j++ {
           // minimal sum for kid j: k*(j+1) + base
           req := kk*int64(j+1) + base
           if req > a[j] {
               return false
           }
       }
       return true
   }
   for lo+1 < hi {
       mid := (lo + hi) >> 1
       if check(mid) {
           lo = mid
       } else {
           hi = mid
       }
   }
   k := lo
   if k == 0 {
       fmt.Fprint(out, 0)
       return
   }
   // compute minimal total and leftover budgets
   kk := int64(k)
   kn := int64(n) * kk
   // minimal sum of sequence 1..kn
   minTotal := kn * (kn + 1) / 2
   // per-kid minimal sum
   base := int64(n) * kk * (kk - 1) / 2
   bLeft := make([]int64, n)
   for j := 0; j < n; j++ {
       // minimal sum for kid j: k*(j+1) + base
       minj := kk*int64(j+1) + base
       bLeft[j] = a[j] - minj
       // rem unused in descending greedy
   }
   // no global cap needed for descending greedy
   var extra int64 = 0
   // greedy assign extra x
   // descending greedy: assign s[i] from i down to 1
   // nextVal is s[i+1], init above maximal m+1
   nextVal := int64(m) + 1
   for i := kn; i >= 1; i-- {
       idx := int((i - 1) % int64(n))
       // max allowed by monotonic: s[i] < nextVal
       maxMon := nextVal - 1
       // max allowed by budget: minimal i + bLeft
       maxBud := i + bLeft[idx]
       // cap by global m constraint
       if maxMon > int64(m) {
           maxMon = int64(m)
       }
       // choose s_i
       var si int64
       if maxBud < maxMon {
           si = maxBud
       } else {
           si = maxMon
       }
       if si < i {
           si = i
       }
       // update
       extra += (si - i)
       bLeft[idx] -= (si - i)
       nextVal = si
   }
   total := minTotal + extra
   fmt.Fprint(out, total)
}
