package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// min64 returns the smaller of two int64 values.
func min64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

// abs64 returns the absolute value of an int64.
func abs64(a int64) int64 {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // transform values
   d := make([]int64, n)
   for i := 0; i < n; i++ {
       d[i] = a[i] - int64(i)
   }
   // candidates: sorted unique d
   v := make([]int64, n)
   copy(v, d)
   sort.Slice(v, func(i, j int) bool { return v[i] < v[j] })
   w := make([]int64, 0, n)
   for i, val := range v {
       if i == 0 || val != v[i-1] {
           w = append(w, val)
       }
   }
   m := len(w)
   // DP arrays
   dpPrev := make([]int64, m)
   dpCur := make([]int64, m)
   // base case for first element
   for j := 0; j < m; j++ {
       dpPrev[j] = abs64(d[0] - w[j])
   }
   // DP transitions
   for i := 1; i < n; i++ {
       // prefix minima
       pref := make([]int64, m)
       pref[0] = dpPrev[0]
       for j := 1; j < m; j++ {
           pref[j] = min64(pref[j-1], dpPrev[j])
       }
       // compute dpCur
       for j := 0; j < m; j++ {
           dpCur[j] = pref[j] + abs64(d[i] - w[j])
       }
       // swap for next iteration
       dpPrev, dpCur = dpCur, dpPrev
   }
   // answer is min over last dpPrev
   ans := dpPrev[0]
   for j := 1; j < m; j++ {
       if dpPrev[j] < ans {
           ans = dpPrev[j]
       }
   }
   fmt.Fprintln(writer, ans)
}
