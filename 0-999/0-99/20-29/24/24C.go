package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   rdr := bufio.NewReader(os.Stdin)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   var n int
   var j int64
   if _, err := fmt.Fscan(rdr, &n, &j); err != nil {
       return
   }
   var x0, y0 int64
   fmt.Fscan(rdr, &x0, &y0)
   Ax := make([]int64, n)
   Ay := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(rdr, &Ax[i], &Ay[i])
   }
   // compute C = sum_{t=0..n-1} (-1)^t * A[t]
   var Cx, Cy int64
   for t := 0; t < n; t++ {
       if t%2 == 0 {
           Cx += Ax[t]
           Cy += Ay[t]
       } else {
           Cx -= Ax[t]
           Cy -= Ay[t]
       }
   }
   // full cycles
   q := j / int64(n)
   r := int(j % int64(n))
   // base point after q*n steps
   var mx, my int64
   if q%2 == 0 {
       mx, my = x0, y0
   } else {
       mx = -x0 + 2*Cx
       my = -y0 + 2*Cy
   }
   // apply remaining r steps
   for t := 0; t < r; t++ {
       // reflect mx,my across A[t]
       mx = -mx + 2*Ax[t]
       my = -my + 2*Ay[t]
   }
   fmt.Fprintf(w, "%d %d", mx, my)
}
