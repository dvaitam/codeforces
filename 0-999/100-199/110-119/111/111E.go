package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m, x1, y1, x2, y2 int
   fmt.Fscan(reader, &n, &m)
   fmt.Fscan(reader, &x1, &y1)
   fmt.Fscan(reader, &x2, &y2)
   // Simple path: go horizontally then vertically
   var path [][2]int
   // horizontal move
   if y1 <= y2 {
       for y := y1; y <= y2; y++ {
           path = append(path, [2]int{x1, y})
       }
   } else {
       for y := y1; y >= y2; y-- {
           path = append(path, [2]int{x1, y})
       }
   }
   // vertical move (skip first cell at (x1,y2))
   if x1 < x2 {
       for x := x1 + 1; x <= x2; x++ {
           path = append(path, [2]int{x, y2})
       }
   } else if x1 > x2 {
       for x := x1 - 1; x >= x2; x-- {
           path = append(path, [2]int{x, y2})
       }
   }
   // output
   fmt.Fprintln(writer, len(path))
   for _, p := range path {
       fmt.Fprintln(writer, p[0], p[1])
   }
}
