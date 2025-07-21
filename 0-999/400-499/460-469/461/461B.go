package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func add(a, b int) int {
   a += b
   if a >= mod {
       a -= mod
   }
   return a
}

func mul(a, b int) int {
   return int((int64(a) * int64(b)) % mod)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   parents := make([]int, n)
   // parents[0] unused
   for i := 1; i < n; i++ {
       var p int
       fmt.Fscan(reader, &p)
       parents[i] = p
   }
   color := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &color[i])
   }
   // build children
   children := make([][]int, n)
   for i := 1; i < n; i++ {
       p := parents[i]
       children[p] = append(children[p], i)
   }
   // dp0_0, dp0_1, dp0_2 for each node
   dp0 := make([][3]int, n)
   // init base dp with node itself before children
   for u := 0; u < n; u++ {
       if color[u] == 1 {
           dp0[u][0] = 0
           dp0[u][1] = 1
           dp0[u][2] = 0
       } else {
           dp0[u][0] = 1
           dp0[u][1] = 0
           dp0[u][2] = 0
       }
   }
   // post-order DP from n-1 down to 0 (parents have smaller index)
   for u := n - 1; u >= 0; u-- {
       // combine children
       for _, v := range children[u] {
           // t = dp0[v]
           t0 := dp0[v][0]
           t1 := dp0[v][1]
           t2 := dp0[v][2]
           cur := dp0[u]
           // reset
           dp0[u][0], dp0[u][1], dp0[u][2] = 0, 0, 0
           // for each state j of u before this child
           for j := 0; j < 3; j++ {
               cj := cur[j]
               if cj == 0 {
                   continue
               }
               // option: cut edge u-v, child must have exactly one black on its main component -> t1
               dp0[u][j] = add(dp0[u][j], mul(cj, t1))
               // option: merge child into u's main component
               // child main comp has k blacks: contribute to j+k
               if t0 != 0 {
                   nj := j
                   dp0[u][nj] = add(dp0[u][nj], mul(cj, t0))
               }
               if t1 != 0 {
                   nj := j + 1
                   if nj > 2 {
                       nj = 2
                   }
                   dp0[u][nj] = add(dp0[u][nj], mul(cj, t1))
               }
               if t2 != 0 {
                   nj := 2
                   dp0[u][nj] = add(dp0[u][nj], mul(cj, t2))
               }
           }
       }
   }
   // answer: dp0[0][1]
   fmt.Fprintln(writer, dp0[0][1])
}
