package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var R, G, B int
   if _, err := fmt.Fscan(reader, &R, &G, &B); err != nil {
       return
   }
   r := make([]int, R)
   g := make([]int, G)
   b := make([]int, B)
   for i := 0; i < R; i++ {
       fmt.Fscan(reader, &r[i])
   }
   for i := 0; i < G; i++ {
       fmt.Fscan(reader, &g[i])
   }
   for i := 0; i < B; i++ {
       fmt.Fscan(reader, &b[i])
   }
   // sort in decreasing order
   sort.Slice(r, func(i, j int) bool { return r[i] > r[j] })
   sort.Slice(g, func(i, j int) bool { return g[i] > g[j] })
   sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
   // 1-based
   r2 := make([]int, R+1)
   g2 := make([]int, G+1)
   b2 := make([]int, B+1)
   for i := 1; i <= R; i++ {
       r2[i] = r[i-1]
   }
   for i := 1; i <= G; i++ {
       g2[i] = g[i-1]
   }
   for i := 1; i <= B; i++ {
       b2[i] = b[i-1]
   }
   // dpPrev and dpCurr with dimensions (R+1)x(G+1)
   dpPrev := make([][]int, R+1)
   dpCurr := make([][]int, R+1)
   for i := 0; i <= R; i++ {
       dpPrev[i] = make([]int, G+1)
       dpCurr[i] = make([]int, G+1)
   }
   // dp for red-green only
   for i := 1; i <= R; i++ {
       for j := 1; j <= G; j++ {
           a := dpPrev[i-1][j-1] + r2[i]*g2[j]
           b0 := dpPrev[i-1][j]
           c := dpPrev[i][j-1]
           // take max of a, b0, c
           m := a
           if b0 > m {
               m = b0
           }
           if c > m {
               m = c
           }
           dpPrev[i][j] = m
       }
   }
   // add blue sticks one by one
   for k := 1; k <= B; k++ {
       // start with skipping b[k]
       for i := 0; i <= R; i++ {
           copy(dpCurr[i], dpPrev[i])
       }
       // red-blue
       for i := 1; i <= R; i++ {
           vi := r2[i] * b2[k]
           for j := 0; j <= G; j++ {
               v := dpPrev[i-1][j] + vi
               if v > dpCurr[i][j] {
                   dpCurr[i][j] = v
               }
           }
       }
       // green-blue
       for i := 0; i <= R; i++ {
           for j := 1; j <= G; j++ {
               v := dpPrev[i][j-1] + g2[j]*b2[k]
               if v > dpCurr[i][j] {
                   dpCurr[i][j] = v
               }
           }
       }
       // swap
       dpPrev, dpCurr = dpCurr, dpPrev
   }
   // answer
   fmt.Fprintln(writer, dpPrev[R][G])
}
