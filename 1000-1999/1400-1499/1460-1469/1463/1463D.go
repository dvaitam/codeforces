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
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(in, &n)
       b := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &b[i])
       }
       // build complement array a
       used := make([]bool, 2*n+1)
       for i := 0; i < n; i++ {
           if b[i] >= 1 && b[i] <= 2*n {
               used[b[i]] = true
           }
       }
       a := make([]int, 0, n)
       for v := 1; v <= 2*n; v++ {
           if !used[v] {
               a = append(a, v)
           }
       }
       // fmin_count: max matches for b as mins (a > b)
       fmin := 0
       pa, pb := 0, 0
       for pa < n && pb < n {
           if a[pa] > b[pb] {
               fmin++
               pa++
               pb++
           } else {
               pa++
           }
       }
       // fmax_count: max matches for b as maxs (a < b)
       fmax := 0
       pa = n - 1
       pb = n - 1
       for pa >= 0 && pb >= 0 {
           if a[pa] < b[pb] {
               fmax++
               pa--
               pb--
           } else {
               pa--
           }
       }
       // x in [0..fmin] and n-x <= fmax => x >= n-fmax
       low := n - fmax
       if low < 0 {
           low = 0
       }
       high := fmin
       if high > n {
           high = n
       }
       ans := 0
       if high >= low {
           ans = high - low + 1
       }
       fmt.Fprintln(out, ans)
   }
}
