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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // Read initial distance matrix
   d := make([][]int, n)
   for i := 0; i < n; i++ {
       d[i] = make([]int, n)
       for j := 0; j < n; j++ {
           fmt.Fscan(reader, &d[i][j])
       }
   }
   var k int
   fmt.Fscan(reader, &k)
   ans := make([]int64, k)
   // Process each new road
   for idx := 0; idx < k; idx++ {
       var u, v, c int
       fmt.Fscan(reader, &u, &v, &c)
       u--
       v--
       // copy distances to u and v
       du := make([]int, n)
       dv := make([]int, n)
       for i := 0; i < n; i++ {
           du[i] = d[i][u]
           dv[i] = d[i][v]
       }
       // update distances using the new road
       for i := 0; i < n; i++ {
           for j := 0; j < n; j++ {
               // possible better path: i->u->v->j or i->v->u->j
               alt := du[i] + c + dv[j]
               if dv[i]+c+du[j] < alt {
                   alt = dv[i] + c + du[j]
               }
               if alt < d[i][j] {
                   d[i][j] = alt
               }
           }
       }
       // compute sum of distances for all unordered pairs
       var sum int64
       for i := 0; i < n; i++ {
           for j := i + 1; j < n; j++ {
               sum += int64(d[i][j])
           }
       }
       ans[idx] = sum
   }
   // Output results
   for i, v := range ans {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
