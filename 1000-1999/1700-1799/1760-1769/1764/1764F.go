package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var N int
   if _, err := fmt.Fscan(in, &N); err != nil {
       return
   }
   A := make([][]int64, N)
   for i := 0; i < N; i++ {
       A[i] = make([]int64, N)
       for j := 0; j <= i; j++ {
           var x int64
           fmt.Fscan(in, &x)
           A[i][j] = x
           A[j][i] = x
       }
   }
   sz := make([]int64, N)
   alive := make([]bool, N)
   for i := 0; i < N; i++ {
       sz[i] = 1
       alive[i] = true
   }
   for z := 0; z < N-1; z++ {
       v := -1
       for i := 0; i < N; i++ {
           if !alive[i] {
               continue
           }
           if v < 0 || A[i][i] > A[v][v] {
               v = i
           }
       }
       alive[v] = false
       u := -1
       for i := 0; i < N; i++ {
           if !alive[i] {
               continue
           }
           if u < 0 || A[v][i] > A[v][u] {
               u = i
           }
       }
       w := (A[u][u] - A[v][u]) / sz[v]
       sz[u] += sz[v]
       fmt.Fprintf(out, "%d %d %d\n", v+1, u+1, w)
   }
}
