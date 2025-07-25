package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   total := 2 * n
   a := make([]int64, total)
   b := make([]int64, total)
   for i := 0; i < total; i++ {
       fmt.Fscan(reader, &a[i], &b[i])
   }
   // binary search minimal k
   var lo, hi int64 = 0, 2000000000
   ans := int64(-1)
   for lo <= hi {
       mid := (lo + hi) / 2
       if feasible(mid, n, a, b) {
           ans = mid
           hi = mid - 1
       } else {
           lo = mid + 1
       }
   }
   if ans < 0 {
       fmt.Println(-1)
   } else {
       // minimal number of exchanges: assume direct swaps
       fmt.Printf("%d %d\n", ans, n)
   }
}

// feasible checks if at time k, initial stocks can be matched directly to targets
func feasible(k int64, n int, a, b []int64) bool {
   initP := make([]int64, n)
   targP := make([]int64, n)
   for i := 0; i < n; i++ {
       initP[i] = a[i]*k + b[i]
   }
   for j := 0; j < n; j++ {
       targP[j] = a[n+j]*k + b[n+j]
   }
   sort.Slice(initP, func(i, j int) bool { return initP[i] < initP[j] })
   sort.Slice(targP, func(i, j int) bool { return targP[i] < targP[j] })
   for i := 0; i < n; i++ {
       if initP[i] < targP[i] {
           return false
       }
   }
   return true
}
