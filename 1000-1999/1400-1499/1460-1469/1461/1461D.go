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

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n, q int
       fmt.Fscan(reader, &n, &q)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       sort.Ints(a)
       // prefix sums
       pref := make([]int64, n+1)
       for i := 0; i < n; i++ {
           pref[i+1] = pref[i] + int64(a[i])
       }
       // collect possible sums
       sums := make(map[int64]bool)
       type seg struct{ l, r int }
       stack := make([]seg, 0, 2*n)
       stack = append(stack, seg{0, n - 1})
       for len(stack) > 0 {
           s := stack[len(stack)-1]
           stack = stack[:len(stack)-1]
           l, r := s.l, s.r
           sum := pref[r+1] - pref[l]
           if sums[sum] {
               continue
           }
           sums[sum] = true
           if a[l] == a[r] {
               continue
           }
           mid := (a[l] + a[r]) / 2
           // find last index <= mid
           cnt := r - l + 1
           idx := sort.Search(cnt, func(i int) bool { return a[l+i] > mid }) - 1
           pivot := l + idx
           // left segment [l..pivot]
           if pivot >= l {
               stack = append(stack, seg{l, pivot})
           }
           // right segment [pivot+1..r]
           if pivot+1 <= r {
               stack = append(stack, seg{pivot + 1, r})
           }
       }
       // answer queries
       for i := 0; i < q; i++ {
           var s int64
           fmt.Fscan(reader, &s)
           if sums[s] {
               fmt.Fprintln(writer, "Yes")
           } else {
               fmt.Fprintln(writer, "No")
           }
       }
   }
}
