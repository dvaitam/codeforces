package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, A, B int
   if _, err := fmt.Fscan(reader, &n, &m, &A, &B); err != nil {
       return
   }
   A--
   B--
   if A < B {
       A, B = B, A
   }
   g := make([][]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       g[u] = append(g[u], v)
       g[v] = append(g[v], u)
   }
   p := make([]float64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   id := make([][]int, n)
   for i := 0; i < n; i++ {
       id[i] = make([]int, n)
   }
   N := 0
   for i := 0; i < n; i++ {
       for j := 0; j <= i; j++ {
           id[i][j] = N
           N++
       }
   }
   a := make([][]float64, N)
   for i := 0; i < N; i++ {
       a[i] = make([]float64, N)
   }
   edge := func(ea, b1, c, d int, e float64) {
       if ea < b1 {
           ea, b1 = b1, ea
       }
       if c < d {
           c, d = d, c
       }
       i1 := id[ea][b1]
       i2 := id[c][d]
       a[i1][i2] += e
   }
   for i := 0; i < n; i++ {
       for j := 0; j <= i; j++ {
           if i == j {
               edge(i, j, i, j, 1)
               continue
           }
           edge(i, j, i, j, p[i]*p[j])
           giLen := float64(len(g[i]))
           gjLen := float64(len(g[j]))
           for _, qi := range g[i] {
               edge(i, j, qi, j, (1-p[i])/giLen*p[j])
               for _, qj := range g[j] {
                   edge(i, j, qi, qj, (1-p[i])/giLen*(1-p[j])/gjLen)
               }
           }
           for _, qj := range g[j] {
               edge(i, j, i, qj, (1-p[j])/gjLen*p[i])
           }
       }
   }
   for i := 0; i < n; i++ {
       for j := 0; j < i; j++ {
           idx := id[i][j]
           a[idx][idx] -= 1
       }
   }
   l := make([][]float64, N)
   u := make([][]float64, N)
   for i := 0; i < N; i++ {
       l[i] = make([]float64, N)
       u[i] = make([]float64, N)
   }
   for j := 0; j < N; j++ {
       u[0][j] = a[0][j]
   }
   for j := 1; j < N; j++ {
       l[j][0] = a[j][0] / u[0][0]
   }
   for i := 1; i < N; i++ {
       for j := i; j < N; j++ {
           sum := a[i][j]
           for k := 0; k < i; k++ {
               sum -= l[i][k] * u[k][j]
           }
           u[i][j] = sum
       }
       for j := i + 1; j < N; j++ {
           sum := a[j][i]
           for k := 0; k < i; k++ {
               sum -= l[j][k] * u[k][i]
           }
           l[j][i] = sum / u[i][i]
       }
   }
   for i := 0; i < N; i++ {
       l[i][i] = 1
   }
   x := make([]float64, N)
   y := make([]float64, N)
   bvec := make([]float64, N)
   res := make([]float64, n)
   idxAB := id[A][B]
   for d := 0; d < n; d++ {
       for i := 0; i < N; i++ {
           bvec[i] = 0
       }
       bvec[id[d][d]] = 1
       for i := 0; i < N; i++ {
           sum := bvec[i]
           for j := 0; j < i; j++ {
               sum -= l[i][j] * y[j]
           }
           y[i] = sum
       }
       for i := N - 1; i >= 0; i-- {
           sum := y[i]
           for j := i + 1; j < N; j++ {
               sum -= u[i][j] * x[j]
           }
           x[i] = sum / u[i][i]
       }
       res[d] = x[idxAB]
   }
   for i, v := range res {
       if i > 0 {
           fmt.Print(" ")
       }
       fmt.Printf("%.8f", v)
   }
   fmt.Println()
}
