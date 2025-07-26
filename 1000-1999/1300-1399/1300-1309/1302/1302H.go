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

   var n, q int
   if _, err := fmt.Fscan(reader, &n, &q); err != nil {
       return
   }
   a := make([]uint64, n+1)
   for i := 1; i <= n; i++ {
       var x uint64
       fmt.Fscan(reader, &x)
       // shift by 1 to avoid zero confusion
       a[i] = x + 1
   }

   // rolling hash parameters
   const base1 = 91138233
   const base2 = 97266353
   h1 := make([]uint64, n+1)
   h2 := make([]uint64, n+1)
   p1 := make([]uint64, n+1)
   p2 := make([]uint64, n+1)
   p1[0], p2[0] = 1, 1
   for i := 1; i <= n; i++ {
       p1[i] = p1[i-1] * base1
       p2[i] = p2[i-1] * base2
       h1[i] = h1[i-1]*base1 + a[i]
       h2[i] = h2[i-1]*base2 + a[i]
   }

   for qi := 0; qi < q; qi++ {
       var length, p, qpos int
       fmt.Fscan(reader, &length, &p, &qpos)
       l1, r1 := p, p+length-1
       l2, r2 := qpos, qpos+length-1

       // compute hash for [l, r]
       x1 := h1[r1] - h1[l1-1]*p1[r1-l1+1]
       y1 := h1[r2] - h1[l2-1]*p1[r2-l2+1]
       if x1 != y1 {
           fmt.Fprintln(writer, "No")
           continue
       }
       x2 := h2[r1] - h2[l1-1]*p2[r1-l1+1]
       y2 := h2[r2] - h2[l2-1]*p2[r2-l2+1]
       if x2 != y2 {
           fmt.Fprintln(writer, "No")
       } else {
           fmt.Fprintln(writer, "Yes")
       }
   }
}
