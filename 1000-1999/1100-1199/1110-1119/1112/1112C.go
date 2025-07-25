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

   var m, k, n, s int
   if _, err := fmt.Fscan(reader, &m, &k, &n, &s); err != nil {
       return
   }
   a := make([]int, m+1)
   for i := 1; i <= m; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // needed counts
   const MAXV = 500000
   needed := make([]int, MAXV+1)
   b := make([]int, s)
   for i := 0; i < s; i++ {
       fmt.Fscan(reader, &b[i])
       needed[b[i]]++
   }
   // sliding window to cover all b
   cur := make([]int, MAXV+1)
   matched := 0
   l := 1
   r := 0
   for rr := 1; rr <= m; rr++ {
       cur[a[rr]]++
       if cur[a[rr]] <= needed[a[rr]] {
           matched++
       }
       if matched == s {
           // shrink left
           for l <= rr && cur[a[l]] > needed[a[l]] {
               cur[a[l]]--
               l++
           }
           r = rr
           break
       }
   }
   if matched < s || (n > 1 && r > n*k) {
       fmt.Fprintln(writer, -1)
       return
   }
   // positions to remove
   totalRemove := m - n*k
   removes := make([]int, 0, totalRemove)
   // copy needed to tmpNeed
   tmpNeed := make([]int, MAXV+1)
   copy(tmpNeed, needed)
   // remove excess inside [l..r] to leave exactly k
   needRemIn := (r - l + 1) - k
   if needRemIn < 0 {
       needRemIn = 0
   }
   for i := l; i <= r && needRemIn > 0; i++ {
       ai := a[i]
       if tmpNeed[ai] > 0 {
           tmpNeed[ai]--
       } else {
           removes = append(removes, i)
           needRemIn--
       }
   }
   // remove to align block start
   remBefore := (l - 1) % k
   for i := 0; i < remBefore; i++ {
       removes = append(removes, (l-1)-i)
   }
   // remove remaining from after r
   for i := r + 1; i <= m && len(removes) < totalRemove; i++ {
       removes = append(removes, i)
   }
   // output
   if len(removes) > totalRemove {
       removes = removes[:totalRemove]
   }
   sort.Ints(removes)
   fmt.Fprintln(writer, len(removes))
   for i, pos := range removes {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, pos)
   }
   if len(removes) > 0 {
       writer.WriteByte('\n')
   }
}
