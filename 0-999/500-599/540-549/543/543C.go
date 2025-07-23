package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   s := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &s[i])
   }
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &a[i][j])
       }
   }
   // features: position j and letter c, total m*26
   cols := m * 26
   // cost matrix: rows n, cols m*26
   cost := make([][]int, n)
   for i := 0; i < n; i++ {
       cost[i] = make([]int, cols)
   }
   // build costs
   for j := 0; j < m; j++ {
       for c := 0; c < 26; c++ {
           fid := j*26 + c
           // preSum: cost to change others so they != c
           sum := 0
           cc := byte('a' + c)
           for i := 0; i < n; i++ {
               if s[i][j] == cc {
                   sum += a[i][j]
               }
           }
           for i := 0; i < n; i++ {
               // cost for i to have this feature
               cur := sum
               if s[i][j] == cc {
                   cur -= a[i][j]
               } else {
                   cur += a[i][j]
               }
               cost[i][fid] = cur
           }
       }
   }
   // Hungarian for rectangular n x cols, n <= cols
   const INF = 1 << 60
   // potentials
   u := make([]int, n+1)
   v := make([]int, cols+1)
   p := make([]int, cols+1)
   way := make([]int, cols+1)
   for i := 1; i <= n; i++ {
       p[0] = i
       j0 := 0
       minv := make([]int, cols+1)
       used := make([]bool, cols+1)
       for j := 0; j <= cols; j++ {
           minv[j] = INF
       }
       for {
           used[j0] = true
           i0 := p[j0]
           delta := INF
           j1 := 0
           for j := 1; j <= cols; j++ {
               if !used[j] {
                   cur := cost[i0-1][j-1] - u[i0] - v[j]
                   if cur < minv[j] {
                       minv[j] = cur
                       way[j] = j0
                   }
                   if minv[j] < delta {
                       delta = minv[j]
                       j1 = j
                   }
               }
           }
           for j := 0; j <= cols; j++ {
               if used[j] {
                   u[p[j]] += delta
                   v[j] -= delta
               } else {
                   minv[j] -= delta
               }
           }
           j0 = j1
           if p[j0] == 0 {
               break
           }
       }
       for {
           j1 := way[j0]
           p[j0] = p[j1]
           j0 = j1
           if j0 == 0 {
               break
           }
       }
   }
   // p[j] = row assigned to column j
   assignment := make([]int, n)
   for j := 1; j <= cols; j++ {
       if p[j] > 0 {
           assignment[p[j]-1] = j - 1
       }
   }
   // compute answer
   ans := 0
   for i := 0; i < n; i++ {
       ans += cost[i][assignment[i]]
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}
