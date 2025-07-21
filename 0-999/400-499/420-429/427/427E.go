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

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // Precompute sums for left jumps of m
   sL := make([]int64, n)
   for i := 0; i < n; i++ {
       if i >= m {
           sL[i] = a[i] + sL[i-m]
       } else {
           sL[i] = a[i]
       }
   }
   // Precompute sums for right jumps of m
   sR := make([]int64, n)
   for i := n - 1; i >= 0; i-- {
       if i+m < n {
           sR[i] = a[i] + sR[i+m]
       } else {
           sR[i] = a[i]
       }
   }
   var best int64 = -1
   // iterate unique positions
   for i := 0; i < n; {
       j := i + 1
       for j < n && a[j] == a[i] {
           j++
       }
       // block [i..j-1] with value S
       S := a[i]
       // left side: indices [0..i-1]
       L := int64(i)
       var sumL int64
       if i > 0 {
           sumL = sL[i-1]
       }
       TL := (L + int64(m) - 1) / int64(m)
       distL := TL*S - sumL
       // right side: indices [j..n-1]
       Rn := int64(n - j)
       var sumRsel int64
       TR := (Rn + int64(m) - 1) / int64(m)
       if Rn > 0 {
           // start index of first in selection
           // index = (n-1) - (TR-1)*m
           idx0 := int64(n-1) - (TR-1)*int64(m)
           if idx0 < int64(j) {
               // ensure starting at or after j
               // adjust to the first idx >= j with same mod m
               modR := (int64(n-1)) % int64(m)
               rem := int64(j) % int64(m)
               delta := (modR - rem + int64(m)) % int64(m)
               idx0 = int64(j) + delta
           }
           if idx0 < int64(n) {
               sumRsel = sR[idx0]
           }
       }
       distR := sumRsel - TR*S
       total := 2 * (distL + distR)
       if best < 0 || total < best {
           best = total
       }
       i = j
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprint(w, best)
}
