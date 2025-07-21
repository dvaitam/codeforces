package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type segment struct {
   l, r int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   segs := make([]segment, 0, n)
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       if x <= y {
           segs = append(segs, segment{x, y})
       } else {
           segs = append(segs, segment{y, x})
       }
   }
   // sort by right endpoint
   sort.Slice(segs, func(i, j int) bool {
       return segs[i].r < segs[j].r
   })
   nails := make([]int, 0, n)
   hasLast := false
   last := 0
   for _, s := range segs {
       if !hasLast || last < s.l {
           // place a nail at right endpoint
           last = s.r
           nails = append(nails, last)
           hasLast = true
       }
   }
   // output result
   fmt.Fprintln(writer, len(nails))
   for i, v := range nails {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
