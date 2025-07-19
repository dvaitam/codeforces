package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   var k int
   if _, err := fmt.Fscan(os.Stdin, &k); err != nil {
       return
   }
   const n = 95
   edge := make([][]bool, n)
   for i := range edge {
       edge[i] = make([]bool, n)
   }
   draw := func(a, b int) {
       edge[a][b] = true
       edge[b][a] = true
   }
   // Build 2x2 bipartite blocks
   for i := 0; i < 30; i++ {
       for j := 2*i + 2; j < 2*i + 4; j++ {
           for l := 2*i + 4; l < 2*i + 6; l++ {
               draw(j, l)
           }
       }
   }
   // Build path from 64 to 94
   for i := 64; i < 94; i++ {
       draw(i, i+1)
   }
   // Connect start and end anchors
   draw(0, 2)
   draw(0, 3)
   draw(94, 1)
   // Connect blocks to path based on bits of k
   for i := 0; i < 30; i++ {
       if k&(1<<i) != 0 {
           draw(2*i+2, 64+i+1)
       }
   }
   // Output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, n)
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if edge[i][j] {
               writer.WriteByte('Y')
           } else {
               writer.WriteByte('N')
           }
       }
       writer.WriteByte('\n')
   }
}
