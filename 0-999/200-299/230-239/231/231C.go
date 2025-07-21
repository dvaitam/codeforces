package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   pref := make([]int64, n+1)
   for i := 0; i < n; i++ {
       pref[i+1] = pref[i] + a[i]
   }
   l := 0
   bestCount := 1
   bestValue := a[0]
   for r := 0; r < n; r++ {
       // shrink window until cost <= k
       for l <= r && int64(r-l+1)*a[r] - (pref[r+1] - pref[l]) > k {
           l++
       }
       cnt := r - l + 1
       if cnt > bestCount || (cnt == bestCount && a[r] < bestValue) {
           bestCount = cnt
           bestValue = a[r]
       }
   }
   fmt.Fprint(writer, bestCount, " ", bestValue)
}
