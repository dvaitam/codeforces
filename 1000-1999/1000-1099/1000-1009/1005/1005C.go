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

   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   sort.Ints(a)
   // Precompute powers of two up to 2^30
   pow2 := make([]int, 31)
   for i := 1; i <= 30; i++ {
       pow2[i] = 1 << i
   }
   ans := 0
   for i := 0; i < n; i++ {
       for j := 1; j <= 30; j++ {
           if pow2[j] <= a[i] {
               continue
           }
           target := pow2[j] - a[i]
           idx := sort.SearchInts(a, target)
           if idx < n && a[idx] == target {
               if idx != i || (idx > 0 && a[idx-1] == target) || (idx+1 < n && a[idx+1] == target) {
                   ans++
                   break
               }
           }
       }
   }
   fmt.Fprint(out, n-ans)
}
