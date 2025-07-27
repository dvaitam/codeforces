package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   l := make([]int, m)
   r := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &l[i], &r[i])
   }
   // number of window positions
   w := n - k + 1
   // f[p][i] = overlap length of participant p with window starting at i (1-based)
   f := make([][]int, m)
   for p := 0; p < m; p++ {
       f[p] = make([]int, w+1)
       for i := 1; i <= w; i++ {
           // window [i, i+k-1]
           li, ri := l[p], r[p]
           wi, wj := i, i+k-1
           // overlap is [max(li,wi), min(ri, wj)]
           a := li
           if wi > a {
               a = wi
           }
           b := ri
           if wj < b {
               b = wj
           }
           if a <= b {
               f[p][i] = b - a + 1
           } else {
               f[p][i] = 0
           }
       }
   }
   // base[i] = sum of f[p][i]
   base := make([]int, w+2)
   for i := 1; i <= w; i++ {
       s := 0
       for p := 0; p < m; p++ {
           s += f[p][i]
       }
       base[i] = s
   }
   // best left and right base
   bestLeft := make([]int, w+2)
   best := 0
   for i := 1; i <= w; i++ {
       // windows ending before i: j + k -1 < i => j < i-k+1
       if i-k >= 1 {
           if base[i-k] > bestLeft[i-1] {
               bestLeft[i] = base[i-k]
           } else {
               bestLeft[i] = bestLeft[i-1]
           }
       } else {
           bestLeft[i] = bestLeft[i-1]
       }
   }
   bestRight := make([]int, w+2)
   for i := w; i >= 1; i-- {
       // windows starting after i: j > i+k-1
       if i+k <= w {
           if base[i+k] > bestRight[i+1] {
               bestRight[i] = base[i+k]
           } else {
               bestRight[i] = bestRight[i+1]
           }
       } else {
           bestRight[i] = bestRight[i+1]
       }
   }
   // evaluate best
   for i := 1; i <= w; i++ {
       // using disjoint
       if bestLeft[i] > 0 {
           tmp := base[i] + bestLeft[i]
           if tmp > best {
               best = tmp
           }
       }
       if bestRight[i] > 0 {
           tmp := base[i] + bestRight[i]
           if tmp > best {
               best = tmp
           }
       }
       // same window
       if base[i] > best {
           best = base[i]
       }
       // overlapping windows
       // j from max(1, i-k+1) to min(w, i+k-1)
       lo := i - k + 1
       if lo < 1 {
           lo = 1
       }
       hi := i + k - 1
       if hi > w {
           hi = w
       }
       for j := lo; j <= hi; j++ {
           // skip disjoint and same handled
           // compute sum min
           mn := 0
           for p := 0; p < m; p++ {
               x := f[p][i]
               y := f[p][j]
               if x < y {
                   mn += x
               } else {
                   mn += y
               }
           }
           tmp := base[i] + base[j] - mn
           if tmp > best {
               best = tmp
           }
       }
   }
   fmt.Fprintln(writer, best)
}
