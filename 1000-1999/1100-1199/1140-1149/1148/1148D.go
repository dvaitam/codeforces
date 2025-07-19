package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type segment struct {
   l, r, id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   b := make([]int, n)
   cnt1, cnt2 := 0, 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i], &b[i])
       if a[i] < b[i] {
           cnt1++
       } else {
           cnt2++
       }
   }
   if cnt1 < cnt2 {
       twoN := 2*n
       for i := 0; i < n; i++ {
           a[i] = twoN - a[i] + 1
           b[i] = twoN - b[i] + 1
       }
   }
   var segs []segment
   segs = make([]segment, 0, n)
   for i := 0; i < n; i++ {
       if a[i] < b[i] {
           segs = append(segs, segment{l: a[i], r: b[i], id: i + 1})
       }
   }
   k := len(segs)
   // find segment with maximum r
   idxMax := 0
   for i := 1; i < k; i++ {
       if segs[i].r > segs[idxMax].r {
           idxMax = i
       }
   }
   // output count
   fmt.Fprintln(writer, k)
   if k == 0 {
       return
   }
   // print the segment with max r first
   fmt.Fprint(writer, segs[idxMax].id)
   // remove it from list
   segs = append(segs[:idxMax], segs[idxMax+1:]...)
   if len(segs) > 0 {
       // sort by l ascending
       sort.Slice(segs, func(i, j int) bool {
           return segs[i].l < segs[j].l
       })
       // print remaining in reverse order
       for i := len(segs) - 1; i >= 0; i-- {
           fmt.Fprint(writer, " ", segs[i].id)
       }
   }
   fmt.Fprintln(writer)
}
