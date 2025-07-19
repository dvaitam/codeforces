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

   var xs, ys, n int
   fmt.Fscan(reader, &xs, &ys)
   fmt.Fscan(reader, &n)
   x := make([]int, n)
   y := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x[i], &y[i])
   }
   // distance from start to each object
   d := make([]int, n)
   for i := 0; i < n; i++ {
       dx := x[i] - xs
       dy := y[i] - ys
       d[i] = dx*dx + dy*dy
   }
   // distance between objects
   a := make([][]int, n)
   for i := range a {
       a[i] = make([]int, n)
   }
   for i := 0; i < n; i++ {
       for j := 0; j < i; j++ {
           dx := x[i] - x[j]
           dy := y[i] - y[j]
           a[i][j] = dx*dx + dy*dy
           a[j][i] = a[i][j]
       }
   }
   size := 1 << n
   const INF = 1 << 60
   f := make([]int, size)
   g := make([]int, size)
   for i := range f {
       f[i] = INF
   }
   f[0] = 0
   // DP over subsets
   for mask := 0; mask < size; mask++ {
       if f[mask] == INF {
           continue
       }
       // find first unvisited object
       for i := 0; i < n; i++ {
           bit := 1 << i
           if mask&bit == 0 {
               // go alone and return
               m1 := mask | bit
               cost1 := f[mask] + 2*d[i]
               if cost1 < f[m1] {
                   f[m1] = cost1
                   g[m1] = mask
               }
               // go with one more object
               for j := i + 1; j < n; j++ {
                   bitj := 1 << j
                   if mask&bitj == 0 {
                       m2 := mask | bit | bitj
                       cost2 := f[mask] + d[i] + d[j] + a[i][j]
                       if cost2 < f[m2] {
                           f[m2] = cost2
                           g[m2] = mask
                       }
                   }
               }
               break
           }
       }
   }
   full := (1 << n) - 1
   // output result
   fmt.Fprintln(writer, f[full])
   // reconstruct path
   mask := full
   var path []int
   for mask > 0 {
       path = append(path, 0)
       prev := g[mask]
       diff := mask ^ prev
       for i := 0; i < n; i++ {
           if diff&(1<<i) != 0 {
               path = append(path, i+1)
           }
       }
       mask = prev
   }
   path = append(path, 0)
   // print path
   for _, v := range path {
       fmt.Fprint(writer, v, " ")
   }
   fmt.Fprintln(writer)
}
