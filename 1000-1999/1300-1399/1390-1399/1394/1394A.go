package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, d int
   var m int64
   if _, err := fmt.Fscan(reader, &n, &d, &m); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var small, big []int64
   for _, v := range a {
       if v <= m {
           small = append(small, v)
       } else {
           big = append(big, v)
       }
   }
   // sort descending
   sort.Slice(small, func(i, j int) bool { return small[i] > small[j] })
   sort.Slice(big, func(i, j int) bool { return big[i] > big[j] })
   // prefix sums
   smallPref := make([]int64, len(small)+1)
   for i := 0; i < len(small); i++ {
       smallPref[i+1] = smallPref[i] + small[i]
   }
   bigPref := make([]int64, len(big)+1)
   for i := 0; i < len(big); i++ {
       bigPref[i+1] = bigPref[i] + big[i]
   }
   var ans int64
   // case 0 big
   if len(small) > 0 {
       ans = smallPref[min(len(small), n)]
   }
   // try using i bigs
   for i := 1; i <= len(big); i++ {
       // slots needed: positions for bigs and required muzzled days between them
       slots := (i-1)*(d+1) + 1
       if slots > n {
           break
       }
       // remaining days for small
       rem := n - slots
       cntSmall := min(len(small), rem)
       total := bigPref[i] + smallPref[cntSmall]
       if total > ans {
           ans = total
       }
   }
   fmt.Println(ans)
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}
