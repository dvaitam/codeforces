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

   var a [6]int64
   for i := 0; i < 6; i++ {
       fmt.Fscan(in, &a[i])
   }
   var n int
   fmt.Fscan(in, &n)
   b := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &b[i])
   }
   // Prepare all candidate frets (value, note index)
   type pair struct{ v int64; idx int }
   m := 6 * n
   arr := make([]pair, 0, m)
   for i := 0; i < n; i++ {
       bi := b[i]
       for j := 0; j < 6; j++ {
           // fret = bi - a[j]
           arr = append(arr, pair{bi - a[j], i})
       }
   }
   // Sort by fret value
   sort.Slice(arr, func(i, j int) bool {
       return arr[i].v < arr[j].v
   })
   // Sliding window to cover all notes
   cnt := make([]int, n)
   missing := n
   ans := int64(9e18)
   l := 0
   for r := 0; r < len(arr); r++ {
       p := arr[r]
       if cnt[p.idx] == 0 {
           missing--
       }
       cnt[p.idx]++
       // when all covered, shrink from left
       for missing == 0 && l <= r {
           curDiff := arr[r].v - arr[l].v
           if curDiff < ans {
               ans = curDiff
           }
           // remove left
           li := arr[l].idx
           cnt[li]--
           if cnt[li] == 0 {
               missing++
           }
           l++
       }
   }
   // Output result
   fmt.Fprintln(out, ans)
}
