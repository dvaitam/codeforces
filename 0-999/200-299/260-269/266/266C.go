package main

import (
   "bufio"
   "fmt"
   "os"
)

type op struct { t, i, j int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   aX := make([]int, n)
   aY := make([]int, n)
   for i := 1; i < n; i++ {
       fmt.Fscan(reader, &aX[i], &aY[i])
   }
   var ops []op
   for i := 1; i < n; i++ {
       // fix row position: should be at row i+1
       if aX[i] != i+1 {
           old := aX[i]
           newRow := i + 1
           // apply swap in remaining entries
           for j := i + 1; j < n; j++ {
               if aX[j] == old {
                   aX[j] = newRow
               } else if aX[j] == newRow {
                   aX[j] = old
               }
           }
           ops = append(ops, op{1, old, newRow})
           aX[i] = newRow
       }
       // fix column position: should be <= i
       if aY[i] > i {
           old := aY[i]
           newCol := i
           for j := i + 1; j < n; j++ {
               if aY[j] == old {
                   aY[j] = newCol
               } else if aY[j] == newCol {
                   aY[j] = old
               }
           }
           ops = append(ops, op{2, old, newCol})
           aY[i] = newCol
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, len(ops))
   for _, o := range ops {
       fmt.Fprintf(writer, "%d %d %d\n", o.t, o.i, o.j)
   }
}
