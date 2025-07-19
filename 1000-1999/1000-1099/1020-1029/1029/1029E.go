package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 0x7f7f7f7f

var (
   n   int
   adj [][]int
   f   [][]int
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func min3(a, b, c int) int {
   return min(min(a, b), c)
}

func dp(u, fa int) {
   f[u][0], f[u][1], f[u][2] = 0, 0, 0
   mn := INF
   for _, v := range adj[u] {
       if v == fa {
           continue
       }
       dp(v, u)
       m012 := min3(f[v][0], f[v][1], f[v][2])
       f[u][0] += m012
       m02 := min(f[v][0], f[v][2])
       f[u][1] += m02
       f[u][2] += m02
       diff := f[v][0] - min(f[v][0], f[v][2])
       if diff < mn {
           mn = diff
       }
   }
   if fa != 1 {
       f[u][0]++
   }
   f[u][2] += mn
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n)
   adj = make([][]int, n+1)
   for i := 1; i < n; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   f = make([][]int, n+1)
   for i := 0; i <= n; i++ {
       f[i] = make([]int, 3)
   }
   ans := 0
   for _, c := range adj[1] {
       dp(c, 1)
       for _, v := range adj[c] {
           // include all neighbors; f[1] is zero so safe
           ans += min3(f[v][0], f[v][1], f[v][2])
       }
   }
   fmt.Fprintln(writer, ans)
