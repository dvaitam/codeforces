package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, xc, yc int64
   fmt.Fscan(reader, &n, &m)
   fmt.Fscan(reader, &xc, &yc)
   var k int
   fmt.Fscan(reader, &k)
   var total int64
   const INF = 1 << 60
   for i := 0; i < k; i++ {
       var dx, dy int64
       fmt.Fscan(reader, &dx, &dy)
       var t1, t2 int64
       if dx > 0 {
           t1 = (n - xc) / dx
       } else if dx < 0 {
           t1 = (1 - xc) / dx
       } else {
           t1 = INF
       }
       if dy > 0 {
           t2 = (m - yc) / dy
       } else if dy < 0 {
           t2 = (1 - yc) / dy
       } else {
           t2 = INF
       }
       t := t1
       if t2 < t {
           t = t2
       }
       if t < 0 {
           t = 0
       }
       total += t
       xc += t * dx
       yc += t * dy
   }
   fmt.Println(total)
}
