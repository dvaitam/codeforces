package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// cha holds the bowl capacity b, original index idx, and assigned amount c
type cha struct {
   b   int
   idx int
   c   int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, w int
   fmt.Fscan(reader, &n, &w)
   chas := make([]cha, n)
   sumMin := 0
   for i := 0; i < n; i++ {
       var bi int
       fmt.Fscan(reader, &bi)
       min := (bi + 1) / 2
       chas[i] = cha{b: bi, idx: i, c: min}
       sumMin += min
   }
   if sumMin > w {
       fmt.Fprintln(writer, -1)
       return
   }
   rem := w - sumMin
   // Greedily increase assignments for largest bowls
   sort.Slice(chas, func(i, j int) bool {
       return chas[i].b > chas[j].b
   })
   for i := 0; i < n && rem > 0; i++ {
       avail := chas[i].b - chas[i].c
       if avail > rem {
           avail = rem
       }
       chas[i].c += avail
       rem -= avail
   }
   // Restore original order
   sort.Slice(chas, func(i, j int) bool {
       return chas[i].idx < chas[j].idx
   })
   // Output result
   for i, ch := range chas {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, ch.c)
   }
   writer.WriteByte('\n')
}
