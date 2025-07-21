package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   perimeters := make([]int, n)
   heights := make([]int, n)
   for i := 0; i < n; i++ {
       var l, w, h int
       fmt.Fscan(reader, &l, &w, &h)
       perimeters[i] = 2 * (l + w)
       heights[i] = h
   }
   var m int
   fmt.Fscan(reader, &m)
   types := make([]struct{ L, W, P int }, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &types[i].L, &types[i].W, &types[i].P)
   }
   var total int64
   const inf = int64(1) << 62
   for i := 0; i < n; i++ {
       perim := perimeters[i]
       h := heights[i]
       best := inf
       for j := 0; j < m; j++ {
           L := types[j].L
           W := types[j].W
           P := types[j].P
           strips := L / h
           if strips <= 0 {
               continue
           }
           coverage := strips * W
           rolls := (perim + coverage - 1) / coverage
           cost := int64(rolls) * int64(P)
           if cost < best {
               best = cost
           }
       }
       total += best
   }
   fmt.Println(total)
}
