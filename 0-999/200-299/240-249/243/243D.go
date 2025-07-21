package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type entry struct {
   c, t int64
   h    int64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var vx, vy int64
   if _, err := fmt.Fscan(reader, &n, &vx, &vy); err != nil {
       return
   }
   total := n * n
   entries := make([]entry, 0, total)
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           var h int64
           fmt.Fscan(reader, &h)
           // compute line key c and projection t
           ii := int64(i)
           jj := int64(j)
           c := ii*vy - jj*vx
           t := ii*vx + jj*vy
           if h > 0 {
               entries = append(entries, entry{c: c, t: t, h: h})
           }
           // zero-height towers have no cubes
       }
   }
   // sort by c asc, then t desc
   sort.Slice(entries, func(i, j int) bool {
       if entries[i].c != entries[j].c {
           return entries[i].c < entries[j].c
       }
       return entries[i].t > entries[j].t
   })
   var ans int64
   var curC int64 = 1<<63 - 1
   var maxH int64
   for _, e := range entries {
       if e.c != curC {
           curC = e.c
           maxH = 0
       }
       if e.h > maxH {
           ans += e.h - maxH
           maxH = e.h
       }
   }
   fmt.Fprint(writer, ans)
}
