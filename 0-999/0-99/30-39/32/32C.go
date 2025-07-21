package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, s int64
   if _, err := fmt.Fscan(reader, &n, &m, &s); err != nil {
       return
   }
   // For each dimension, compute number of positions in residue classes with maximum size
   // In 1D of length n and step s:
   // There are s residue classes, each with size either q = n/s or q+1 for r = n%s classes
   qx := n / s
   rx := n % s
   var sumX int64
   if rx > 0 {
       // rx classes have size qx+1
       sumX = rx * (qx + 1)
   } else {
       // all s classes have size qx
       sumX = s * qx
   }
   qy := m / s
   ry := m % s
   var sumY int64
   if ry > 0 {
       sumY = ry * (qy + 1)
   } else {
       sumY = s * qy
   }
   // Total starting positions achieving maximum reachable cells
   result := sumX * sumY
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, result)
}
