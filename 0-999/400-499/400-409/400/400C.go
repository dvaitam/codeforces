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

   var n, m, x, y, z int64
   var p int
   if _, err := fmt.Fscan(reader, &n, &m, &x, &y, &z, &p); err != nil {
       return
   }
   r1 := x % 4
   f := y % 2
   r2 := z % 4

   for k := 0; k < p; k++ {
       var i, j int64
       fmt.Fscan(reader, &i, &j)
       h, w := n, m
       // apply clockwise rotations r1 times
       for t := int64(0); t < r1; t++ {
           i, j = j, h - i + 1
           h, w = w, h
       }
       // apply horizontal flip if needed
       if f == 1 {
           j = w - j + 1
       }
       // apply counterclockwise rotations r2 times
       for t := int64(0); t < r2; t++ {
           i, j = w - j + 1, i
           h, w = w, h
       }
       fmt.Fprintln(writer, i, j)
   }
}
