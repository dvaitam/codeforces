package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       var l int
       fmt.Fscan(reader, &n, &l)
       a := make([]int, n+2)
       a[0] = 0
       for i := 1; i <= n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       a[n+1] = l

       i, j := 0, n+1
       v1, v2 := 1.0, 1.0
       pos1, pos2 := 0.0, float64(l)
       time := 0.0

       for i+1 < j {
           nextL := float64(a[i+1])
           nextR := float64(a[j-1])
           d1 := nextL - pos1
           d2 := pos2 - nextR
           t1 := d1 / v1
           t2 := d2 / v2
           if math.Abs(t1-t2) < 1e-12 {
               // both reach flags simultaneously
               time += t1
               pos1 = nextL
               pos2 = nextR
               v1++
               v2++
               i++
               j--
           } else if t1 < t2 {
               time += t1
               pos1 = nextL
               pos2 -= v2 * t1
               v1++
               i++
           } else {
               time += t2
               pos2 = nextR
               pos1 += v1 * t2
               v2++
               j--
           }
       }
       // remaining distance
       if pos2 > pos1 {
           time += (pos2 - pos1) / (v1 + v2)
       }
       fmt.Fprintf(writer, "%.10f\n", time)
   }
}
