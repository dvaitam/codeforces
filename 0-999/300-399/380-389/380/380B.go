package main

import (
   "bufio"
   "fmt"
   "os"
)

// Op represents a type-1 operation applied to a node range at a given level
type Op struct {
   level, l, r int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // group operations by value x
   opsByX := make(map[int][]Op)
   for i := 0; i < m; i++ {
       var tp int
       fmt.Fscan(reader, &tp)
       if tp == 1 {
           var t, l, r, x int
           fmt.Fscan(reader, &t, &l, &r, &x)
           // record this operation for value x
           opsByX[x] = append(opsByX[x], Op{level: t, l: l, r: r})
       } else if tp == 2 {
           var t0, v0 int
           fmt.Fscan(reader, &t0, &v0)
           // count distinct x that affect subtree of (t0, v0)
           cnt := 0
           for x, ops := range opsByX {
               // for each value x, check if any of its operations hit the subtree
               for _, op := range ops {
                   if op.level < t0 {
                       continue
                   }
                   // at level op.level, subtree of (t0,v0) covers positions [v0 .. v0+(op.level-t0)]
                   maxPos := v0 + (op.level - t0)
                   if op.l <= maxPos && op.r >= v0 {
                       cnt++
                       break
                   }
               }
               // avoid unused x warning
               _ = x
           }
           fmt.Fprintln(writer, cnt)
       }
   }
}
