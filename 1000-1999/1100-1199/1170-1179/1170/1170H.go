package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(in, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
       }
       if n == 0 {
           fmt.Fprintln(out, 0)
           fmt.Fprintln(out)
           continue
       }
       sort.Ints(a)
       // build unique values and counts
       vals := make([]int, 0, n)
       cnt := make([]int, 0, n)
       for i := 0; i < n; {
           j := i + 1
           for j < n && a[j] == a[i] {
               j++
           }
           vals = append(vals, a[i])
           cnt = append(cnt, j-i)
           i = j
       }
       k := len(vals)
       pref := make([]int, k)
       for i := 0; i < k; i++ {
           if i == 0 {
               pref[i] = cnt[i]
           } else {
               pref[i] = pref[i-1] + cnt[i]
           }
       }
       bestM := 1
       bestPattern := 2 // 1: start >, 2: start <
       bestJ := 0
       bestVC, bestPC := 1, 0
       // default: one element
       for j := 0; j < k; j++ {
           vAll := pref[j]
           uAll := n - vAll
           // compute even and odd possibilities
           kEven := vAll
           if uAll < kEven {
               kEven = uAll
           }
           mEven := 2 * kEven
           // odd for pattern1
           kOdd1 := 0
           if uAll > 0 {
               kOdd1 = vAll
               if uAll-1 < kOdd1 {
                   kOdd1 = uAll - 1
               }
           }
           mOdd1 := 2*kOdd1 + 1
           // odd for pattern2
           kOdd2 := 0
           if vAll > 0 {
               kOdd2 = uAll
               if vAll-1 < kOdd2 {
                   kOdd2 = vAll - 1
               }
           }
           mOdd2 := 2*kOdd2 + 1
           // pattern1: start with > (peak first)
           m1, v1, p1 := mEven, kEven, kEven
           if mOdd1 > m1 {
               m1 = mOdd1; v1 = kOdd1; p1 = kOdd1 + 1
           }
           if m1 > bestM {
               bestM = m1; bestPattern = 1; bestJ = j; bestVC = v1; bestPC = p1
           }
           // pattern2: start with < (valley first)
           m2, v2, p2 := mEven, kEven, kEven
           if mOdd2 > m2 {
               m2 = mOdd2; v2 = kOdd2 + 1; p2 = kOdd2
           }
           if m2 > bestM {
               bestM = m2; bestPattern = 2; bestJ = j; bestVC = v2; bestPC = p2
           }
       }
       // reconstruct
       if bestM <= 1 {
           fmt.Fprintln(out, 1)
           if n > 0 {
               fmt.Fprintln(out, a[0])
           } else {
               fmt.Fprintln(out)
           }
           continue
       }
       // split at vals[bestJ]
       threshold := vals[bestJ]
       // build valley_list and peak_list
       valleys := make([]int, 0, pref[bestJ])
       peaks := make([]int, 0, n-pref[bestJ])
       for _, x := range a {
           if x <= threshold {
               valleys = append(valleys, x)
           } else {
               peaks = append(peaks, x)
           }
       }
       // select
       vsel := valleys[:bestVC]
       psel := peaks[:bestPC]
       // output
       fmt.Fprintln(out, bestM)
       // interleave
       vi, pi := 0, 0
       for idx := 0; idx < bestM; idx++ {
           if bestPattern == 1 {
               if idx%2 == 0 {
                   out.WriteString(fmt.Sprintf("%d ", psel[pi])); pi++
               } else {
                   out.WriteString(fmt.Sprintf("%d ", vsel[vi])); vi++
               }
           } else {
               if idx%2 == 0 {
                   out.WriteString(fmt.Sprintf("%d ", vsel[vi])); vi++
               } else {
                   out.WriteString(fmt.Sprintf("%d ", psel[pi])); pi++
               }
           }
       }
       out.WriteString("\n")
   }
}
